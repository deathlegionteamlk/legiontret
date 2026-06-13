"""
LegionTret Python Client
A simple Python library for interacting with the LegionTret API.

Usage:
    from legiontret import Client

    client = Client()
    
    # Generate text
    response = client.generate("llama3", "Why is the sky blue?")
    print(response)
    
    # Chat
    response = client.chat("llama3", [
        {"role": "user", "content": "Hello!"}
    ])
    print(response)
    
    # List models
    models = client.list_models()
    for m in models:
        print(m["name"])
"""

import json
import requests
from typing import List, Dict, Optional, Generator, Any


class LegionTretError(Exception):
    """Base exception for LegionTret errors."""
    pass


class ConnectionError(LegionTretError):
    """Raised when unable to connect to the LegionTret server."""
    pass


class ModelNotFoundError(LegionTretError):
    """Raised when a requested model is not found."""
    pass


class Client:
    """Client for the LegionTret API server.
    
    Args:
        host: Host address of the LegionTret server (default: "127.0.0.1")
        port: Port of the LegionTret server (default: 11434)
        timeout: Request timeout in seconds (default: 300)
    
    Example:
        >>> client = Client()
        >>> response = client.generate("gemma3", "What is Python?")
        >>> print(response["response"])
    """
    
    def __init__(self, host: str = "127.0.0.1", port: int = 11434, timeout: int = 300):
        self.host = host
        self.port = port
        self.timeout = timeout
        self.base_url = f"http://{host}:{port}"
        self._session = requests.Session()
    
    def _url(self, path: str) -> str:
        """Build a full URL from a path."""
        return f"{self.base_url}{path}"
    
    def _request(self, method: str, path: str, json_data: Optional[Dict] = None, 
                 stream: bool = False) -> Any:
        """Make a request to the API."""
        try:
            response = self._session.request(
                method=method,
                url=self._url(path),
                json=json_data,
                timeout=self.timeout,
                stream=stream,
            )
            response.raise_for_status()
            return response
        except requests.exceptions.ConnectionError:
            raise ConnectionError(
                f"Cannot connect to LegionTret at {self.base_url}. "
                "Make sure the server is running: legiontret serve"
            )
        except requests.exceptions.HTTPError as e:
            if e.response.status_code == 404:
                raise ModelNotFoundError(f"Model not found: {json_data}")
            raise LegionTretError(f"HTTP error: {e}")
    
    # ─── Model Management ────────────────────────────────────────────────
    
    def list_models(self) -> List[Dict]:
        """List all locally available models.
        
        Returns:
            List of model dictionaries.
        """
        response = self._request("GET", "/api/tags")
        data = response.json()
        return data.get("models", [])
    
    def show_model(self, name: str) -> Dict:
        """Show details for a specific model.
        
        Args:
            name: Name of the model.
        
        Returns:
            Model details dictionary.
        """
        response = self._request("POST", "/api/show", {"name": name})
        return response.json()
    
    def pull_model(self, name: str, stream: bool = True) -> Generator[Dict, None, None]:
        """Download a model.
        
        Args:
            name: Name of the model to download.
            stream: Whether to stream the download progress.
        
        Yields:
            Progress update dictionaries.
        """
        response = self._request("POST", "/api/pull", {"name": name}, stream=stream)
        if stream:
            for line in response.iter_lines():
                if line:
                    yield json.loads(line)
        else:
            yield response.json()
    
    def delete_model(self, name: str) -> Dict:
        """Delete a local model.
        
        Args:
            name: Name of the model to delete.
        
        Returns:
            Response dictionary.
        """
        response = self._request("DELETE", "/api/delete", {"name": name})
        return response.json()
    
    # ─── Generation ──────────────────────────────────────────────────────
    
    def generate(self, model: str, prompt: str, system: str = "",
                 stream: bool = False, options: Optional[Dict] = None,
                 format: str = "") -> Dict:
        """Generate text from a prompt.
        
        Args:
            model: Name of the model to use.
            prompt: The prompt to generate from.
            system: System prompt (optional).
            stream: Whether to stream the response.
            options: Generation options (temperature, top_p, etc.).
            format: Output format (e.g., "json").
        
        Returns:
            Response dictionary with generated text.
        """
        body = {
            "model": model,
            "prompt": prompt,
            "stream": stream,
        }
        if system:
            body["system"] = system
        if options:
            body["options"] = options
        if format:
            body["format"] = format
        
        response = self._request("POST", "/api/generate", body, stream=stream)
        
        if stream:
            full_response = ""
            for line in response.iter_lines():
                if line:
                    chunk = json.loads(line)
                    full_response += chunk.get("response", "")
                    yield chunk
            return
        
        return response.json()
    
    def generate_stream(self, model: str, prompt: str, system: str = "",
                        options: Optional[Dict] = None) -> Generator[str, None, None]:
        """Generate text from a prompt with streaming.
        
        Args:
            model: Name of the model to use.
            prompt: The prompt to generate from.
            system: System prompt (optional).
            options: Generation options.
        
        Yields:
            Generated text chunks.
        """
        body = {
            "model": model,
            "prompt": prompt,
            "stream": True,
        }
        if system:
            body["system"] = system
        if options:
            body["options"] = options
        
        response = self._request("POST", "/api/generate", body, stream=True)
        for line in response.iter_lines():
            if line:
                chunk = json.loads(line)
                text = chunk.get("response", "")
                if text:
                    yield text
    
    # ─── Chat ────────────────────────────────────────────────────────────
    
    def chat(self, model: str, messages: List[Dict[str, str]], 
             stream: bool = False, options: Optional[Dict] = None,
             format: str = "") -> Dict:
        """Send a chat completion request.
        
        Args:
            model: Name of the model to use.
            messages: List of message dictionaries with 'role' and 'content'.
            stream: Whether to stream the response.
            options: Generation options.
            format: Output format.
        
        Returns:
            Response dictionary.
        
        Example:
            >>> response = client.chat("gemma3", [
            ...     {"role": "system", "content": "You are a helpful assistant."},
            ...     {"role": "user", "content": "What is Python?"},
            ... ])
        """
        body = {
            "model": model,
            "messages": messages,
            "stream": stream,
        }
        if options:
            body["options"] = options
        if format:
            body["format"] = format
        
        response = self._request("POST", "/api/chat", body, stream=stream)
        
        if stream:
            for line in response.iter_lines():
                if line:
                    chunk = json.loads(line)
                    yield chunk
            return
        
        return response.json()
    
    def chat_stream(self, model: str, messages: List[Dict[str, str]],
                    options: Optional[Dict] = None) -> Generator[str, None, None]:
        """Send a chat completion request with streaming.
        
        Args:
            model: Name of the model to use.
            messages: List of message dictionaries.
            options: Generation options.
        
        Yields:
            Generated text chunks.
        """
        body = {
            "model": model,
            "messages": messages,
            "stream": True,
        }
        if options:
            body["options"] = options
        
        response = self._request("POST", "/api/chat", body, stream=True)
        for line in response.iter_lines():
            if line:
                chunk = json.loads(line)
                if "message" in chunk:
                    content = chunk["message"].get("content", "")
                    if content:
                        yield content
    
    # ─── Embeddings ──────────────────────────────────────────────────────
    
    def embeddings(self, model: str, prompt: str) -> Dict:
        """Generate embeddings for a prompt.
        
        Args:
            model: Name of the model to use.
            prompt: The text to embed.
        
        Returns:
            Dictionary containing the embedding vector.
        """
        response = self._request("POST", "/api/embeddings", {
            "model": model,
            "prompt": prompt,
        })
        return response.json()
    
    # ─── OpenAI-Compatible Endpoints ─────────────────────────────────────
    
    def openai_models(self) -> Dict:
        """List models using the OpenAI-compatible endpoint."""
        response = self._request("GET", "/api/v1/models")
        return response.json()
    
    def openai_chat(self, model: str, messages: List[Dict[str, str]],
                    temperature: float = 0.7, max_tokens: int = 2048,
                    stream: bool = False) -> Dict:
        """Chat completion using the OpenAI-compatible endpoint.
        
        This is compatible with any tool that supports the OpenAI API format.
        
        Args:
            model: Name of the model to use.
            messages: List of message dictionaries.
            temperature: Sampling temperature.
            max_tokens: Maximum tokens to generate.
            stream: Whether to stream.
        """
        body = {
            "model": model,
            "messages": messages,
            "temperature": temperature,
            "max_tokens": max_tokens,
            "stream": stream,
        }
        response = self._request("POST", "/api/v1/chat/completions", body, stream=stream)
        return response.json()
    
    # ─── System ──────────────────────────────────────────────────────────
    
    def is_running(self) -> bool:
        """Check if the LegionTret server is running."""
        try:
            self._request("GET", "/health")
            return True
        except Exception:
            return False
    
    def version(self) -> str:
        """Get the LegionTret server version."""
        response = self._request("GET", "/api/version")
        return response.json().get("version", "unknown")
    
    def system_info(self) -> Dict:
        """Get system information from the server."""
        response = self._request("GET", "/api/system")
        return response.json()
    
    def search(self, query: str) -> List[Dict]:
        """Search for models in the registry.
        
        Args:
            query: Search query string.
        
        Returns:
            List of matching models.
        """
        response = self._request("GET", f"/api/search?q={query}")
        data = response.json()
        return data.get("models", [])


# Convenience function for quick generation
def generate(model: str, prompt: str, **kwargs) -> str:
    """Quick generation convenience function.
    
    Args:
        model: Model name.
        prompt: Text prompt.
        **kwargs: Additional options passed to Client.generate().
    
    Returns:
        Generated text string.
    """
    client = Client()
    response = client.generate(model, prompt, **kwargs)
    if isinstance(response, dict):
        return response.get("response", "")
    return ""


def chat(model: str, messages: List[Dict[str, str]], **kwargs) -> str:
    """Quick chat convenience function.
    
    Args:
        model: Model name.
        messages: List of chat messages.
        **kwargs: Additional options passed to Client.chat().
    
    Returns:
        Assistant's response text.
    """
    client = Client()
    response = client.chat(model, messages, **kwargs)
    if isinstance(response, dict):
        msg = response.get("message", {})
        return msg.get("content", "")
    return ""


if __name__ == "__main__":
    # Demo usage
    print("LegionTret Python Client")
    print("=" * 40)
    
    client = Client()
    
    if client.is_running():
        print(f"Server version: {client.version()}")
        print(f"Available models: {len(client.list_models())}")
    else:
        print("LegionTret server is not running.")
        print("Start it with: legiontret serve")
