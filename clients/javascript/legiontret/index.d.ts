/**
 * LegionTret TypeScript declarations
 */
declare module 'legiontret' {
  interface ClientOptions {
    host?: string;
    port?: number;
    timeout?: number;
  }

  interface GenerateOptions {
    system?: string;
    stream?: boolean;
    options?: GenerateParams;
    format?: string;
  }

  interface GenerateParams {
    num_predict?: number;
    temperature?: number;
    top_p?: number;
    top_k?: number;
    repeat_penalty?: number;
    seed?: number;
    stop?: string[];
  }

  interface ChatMessage {
    role: 'system' | 'user' | 'assistant';
    content: string;
    images?: string[];
  }

  interface ModelInfo {
    name: string;
    model: string;
    modified_at: string;
    size: number;
    digest: string;
    details: {
      format: string;
      family: string;
      families: string[];
      parameter_size: string;
      quantization_level: string;
    };
  }

  class LegionTretError extends Error {
    statusCode?: number;
  }

  class ConnectionError extends LegionTretError {}
  class ModelNotFoundError extends LegionTretError {}

  class Client {
    constructor(options?: ClientOptions);
    listModels(): Promise<ModelInfo[]>;
    showModel(name: string): Promise<any>;
    deleteModel(name: string): Promise<any>;
    generate(model: string, prompt: string, options?: GenerateOptions): Promise<any>;
    generateStream(model: string, prompt: string, options?: GenerateOptions): AsyncGenerator<string>;
    chat(model: string, messages: ChatMessage[], options?: any): Promise<any>;
    chatStream(model: string, messages: ChatMessage[], options?: any): AsyncGenerator<string>;
    embeddings(model: string, prompt: string): Promise<any>;
    openaiModels(): Promise<any>;
    openaiChat(model: string, messages: ChatMessage[], options?: any): Promise<any>;
    isRunning(): Promise<boolean>;
    version(): Promise<string>;
    systemInfo(): Promise<any>;
    search(query: string): Promise<any[]>;
  }

  export { Client, LegionTretError, ConnectionError, ModelNotFoundError };
  export default Client;
}
