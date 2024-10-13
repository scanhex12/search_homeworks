from nltk.stem import SnowballStemmer

stemmer = SnowballStemmer('english')
print(stemmer.stem('playing'))