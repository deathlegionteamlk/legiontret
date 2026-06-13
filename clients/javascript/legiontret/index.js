/**
 * LegionTret JavaScript Client
 * A simple JavaScript/TypeScript library for interacting with the LegionTret API.
 *
 * Usage:
 *   import { Client } from 'legiontret';
 *
 *   const client = new Client();
 *
 *   // Generate text
 *   const response = await client.generate('llama3', 'Why is the sky blue?');
 *   console.log(response);
 *
 *   // Chat
 *   const response = await client.chat('llama3', [
 *     { role: 'user', content: 'Hello!' }
 *   ]);
 *   console.log(response);
 *
 *   // List models
 *   const models = await client.listModels();
 *   models.forEach(m => console.log(m.name));
 */

const DEFAULT_HOST = '127.0.0.1';
const DEFAULT_PORT = 11434;
const DEFAULT_TIMEOUT = 300000; // 5 minutes

/**
 * Custom error types for LegionTret
 */
class LegionTretError extends Error {
  constructor(message, statusCode) {
    super(message);
    this.name = 'LegionTretError';
    this.statusCode = statusCode;
  }
}

class ConnectionError extends LegionTretError {
  constructor(host, port) {
    super(`Cannot connect to LegionTret at ${host}:${port}. Make sure the server is running: legiontret serve`);
    this.name = 'ConnectionError';
  }
}

class ModelNotFoundError extends LegionTretError {
  constructor(modelName) {
    super(`Model not found: ${modelName}`);
    this.name = 'ModelNotFoundError';
  }
}

/**
 * LegionTret API Client
 */
class Client {
  /**
   * Create a new LegionTret client.
   * @param {Object} options - Client options
   * @param {string} [options.host='127.0.0.1'] - Host address
   * @param {number} [options.port=11434] - Port number
   * @param {number} [options.timeout=300000] - Request timeout in ms
   */
  constructor(options = {}) {
    this.host = options.host || DEFAULT_HOST;
    this.port = options.port || DEFAULT_PORT;
    this.timeout = options.timeout || DEFAULT_TIMEOUT;
    this.baseUrl = `http://${this.host}:${this.port}`;
  }

  /**
   * Build a full URL from a path
   * @private
   */
  _url(path) {
    return `${this.baseUrl}${path}`;
  }

  /**
   * Make a request to the API
   * @private
   */
  async _request(method, path, body = null, { stream = false } = {}) {
    const options = {
      method,
      headers: { 'Content-Type': 'application/json' },
      signal: AbortSignal.timeout(this.timeout),
    };

    if (body && method !== 'GET') {
      options.body = JSON.stringify(body);
    }

    try {
      const response = await fetch(this._url(path), options);

      if (!response.ok) {
        if (response.status === 404) {
          throw new ModelNotFoundError(body?.model || 'unknown');
        }
        const errorText = await response.text();
        throw new LegionTretError(`HTTP ${response.status}: ${errorText}`, response.status);
      }

      return response;
    } catch (error) {
      if (error.name === 'TypeError' && error.message.includes('fetch')) {
        throw new ConnectionError(this.host, this.port);
      }
      throw error;
    }
  }

  // ─── Model Management ──────────────────────────────────────────────

  /**
   * List all locally available models.
   * @returns {Promise<Array>} List of model objects
   */
  async listModels() {
    const response = await this._request('GET', '/api/tags');
    const data = await response.json();
    return data.models || [];
  }

  /**
   * Show details for a specific model.
   * @param {string} name - Model name
   * @returns {Promise<Object>} Model details
   */
  async showModel(name) {
    const response = await this._request('POST', '/api/show', { name });
    return response.json();
  }

  /**
   * Delete a local model.
   * @param {string} name - Model name to delete
   * @returns {Promise<Object>} Response
   */
  async deleteModel(name) {
    const response = await this._request('DELETE', '/api/delete', { name });
    return response.json();
  }

  // ─── Generation ────────────────────────────────────────────────────

  /**
   * Generate text from a prompt.
   * @param {string} model - Model name
   * @param {string} prompt - Text prompt
   * @param {Object} [options] - Generation options
   * @param {string} [options.system] - System prompt
   * @param {boolean} [options.stream=false] - Stream the response
   * @param {Object} [options.options] - Generation parameters
   * @param {string} [options.format] - Output format
   * @returns {Promise<Object|AsyncGenerator>} Response or stream
   */
  async generate(model, prompt, options = {}) {
    const body = {
      model,
      prompt,
      stream: options.stream || false,
    };
    if (options.system) body.system = options.system;
    if (options.options) body.options = options.options;
    if (options.format) body.format = options.format;

    const response = await this._request('POST', '/api/generate', body, {
      stream: options.stream,
    });

    if (options.stream) {
      return this._parseStream(response);
    }

    return response.json();
  }

  /**
   * Generate text with streaming (async generator).
   * @param {string} model - Model name
   * @param {string} prompt - Text prompt
   * @param {Object} [options] - Generation options
   * @yields {string} Text chunks
   */
  async *generateStream(model, prompt, options = {}) {
    const body = {
      model,
      prompt,
      stream: true,
    };
    if (options.system) body.system = options.system;
    if (options.options) body.options = options.options;

    const response = await this._request('POST', '/api/generate', body, { stream: true });
    yield* this._parseTextStream(response);
  }

  // ─── Chat ──────────────────────────────────────────────────────────

  /**
   * Send a chat completion request.
   * @param {string} model - Model name
   * @param {Array<Object>} messages - Chat messages
   * @param {Object} [options] - Generation options
   * @returns {Promise<Object|AsyncGenerator>} Response or stream
   */
  async chat(model, messages, options = {}) {
    const body = {
      model,
      messages,
      stream: options.stream || false,
    };
    if (options.options) body.options = options.options;
    if (options.format) body.format = options.format;

    const response = await this._request('POST', '/api/chat', body, {
      stream: options.stream,
    });

    if (options.stream) {
      return this._parseStream(response);
    }

    return response.json();
  }

  /**
   * Chat with streaming (async generator).
   * @param {string} model - Model name
   * @param {Array<Object>} messages - Chat messages
   * @param {Object} [options] - Generation options
   * @yields {string} Text chunks
   */
  async *chatStream(model, messages, options = {}) {
    const body = {
      model,
      messages,
      stream: true,
    };
    if (options.options) body.options = options.options;

    const response = await this._request('POST', '/api/chat', body, { stream: true });
    yield* this._parseTextStream(response);
  }

  // ─── Embeddings ────────────────────────────────────────────────────

  /**
   * Generate embeddings for a prompt.
   * @param {string} model - Model name
   * @param {string} prompt - Text to embed
   * @returns {Promise<Object>} Embedding response
   */
  async embeddings(model, prompt) {
    const response = await this._request('POST', '/api/embeddings', { model, prompt });
    return response.json();
  }

  // ─── OpenAI-Compatible ─────────────────────────────────────────────

  /**
   * List models using OpenAI-compatible endpoint.
   * @returns {Promise<Object>} Model list
   */
  async openaiModels() {
    const response = await this._request('GET', '/api/v1/models');
    return response.json();
  }

  /**
   * Chat completion using OpenAI-compatible endpoint.
   * @param {string} model - Model name
   * @param {Array<Object>} messages - Messages
   * @param {Object} [options] - Options
   * @returns {Promise<Object>} Chat completion
   */
  async openaiChat(model, messages, options = {}) {
    const body = {
      model,
      messages,
      temperature: options.temperature || 0.7,
      max_tokens: options.maxTokens || 2048,
      stream: options.stream || false,
    };

    const response = await this._request('POST', '/api/v1/chat/completions', body, {
      stream: options.stream,
    });
    return response.json();
  }

  // ─── System ────────────────────────────────────────────────────────

  /**
   * Check if the LegionTret server is running.
   * @returns {Promise<boolean>}
   */
  async isRunning() {
    try {
      await this._request('GET', '/health');
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Get the server version.
   * @returns {Promise<string>}
   */
  async version() {
    const response = await this._request('GET', '/api/version');
    const data = await response.json();
    return data.version || 'unknown';
  }

  /**
   * Get system information.
   * @returns {Promise<Object>}
   */
  async systemInfo() {
    const response = await this._request('GET', '/api/system');
    return response.json();
  }

  /**
   * Search for models in the registry.
   * @param {string} query - Search query
   * @returns {Promise<Array>} Matching models
   */
  async search(query) {
    const response = await this._request('GET', `/api/search?q=${encodeURIComponent(query)}`);
    const data = await response.json();
    return data.models || [];
  }

  // ─── Private Helpers ───────────────────────────────────────────────

  /**
   * Parse a streaming response as JSON chunks
   * @private
   */
  async *_parseStream(response) {
    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n');
      buffer = lines.pop() || '';

      for (const line of lines) {
        if (line.trim()) {
          try {
            yield JSON.parse(line);
          } catch {
            // Skip malformed lines
          }
        }
      }
    }
  }

  /**
   * Parse a streaming response and yield text content
   * @private
   */
  async *_parseTextStream(response) {
    for await (const chunk of this._parseStream(response)) {
      if (chunk.response) {
        yield chunk.response;
      } else if (chunk.message?.content) {
        yield chunk.message.content;
      }
    }
  }
}

// Export
module.exports = { Client, LegionTretError, ConnectionError, ModelNotFoundError };
module.exports.default = Client;
