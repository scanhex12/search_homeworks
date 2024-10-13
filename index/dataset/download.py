from datasets import load_dataset

ds = load_dataset("glaiveai/RAG-v1")
for text in ds['train']:
    print(text['question'])
    break