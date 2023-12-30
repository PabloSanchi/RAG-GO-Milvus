# Vectorstore with GO and MILVUS
## Embedding Util - Ollama

In this project, I am exploring advanced AI capabilities using Ollama and Milvus. For detailed information about Ollama's AI technologies, you can visit their [website](https://ollama.ai/).

### Using Mistral Model for Embedding with Ollama

My focus is on utilizing the embedding feature of the Mistral model provided by Ollama. This feature converts text data into numerical vectors of 4096 dimensions.

### About Milvus Datastore

[Milvus](https://milvus.io/docs) is an open-source database specifically tailored for storing and managing embeddings. Embeddings are essentially data transformed into numeric vector formats. This functionality is crucial for enabling semantic search, where we can find data based on its meaning and context.

#### Key Functionalities of this project:
- **Create Collection**: Initiating new data collections.
- **Delete Collection**: Removing existing collections.
- **List Collections**: Displaying all available collections.
- **Upsert Documents**: Adding or updating documents in collections.
- **Search Documents**: Conducting searches within collections to find relevant or similar documents.

### Project Purpose

The primary aim of this project is educational, focusing on applying hexagonal architecture in software development. This architecture facilitates a modular and maintainable code structure.

## TODO
- [ ] Dynamic schema
- [ ] Testing