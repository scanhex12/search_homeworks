import spacy
from flask import Flask, request, jsonify
from nltk.corpus import stopwords
from nltk.stem import SnowballStemmer
import nltk

nltk.download('stopwords')

app = Flask(__name__)

nlp = spacy.load("en_core_web_sm")

stemmer = SnowballStemmer("english")

stop_words = set(stopwords.words('english'))

@app.route("/lemmatize", methods=["POST"])
def lemmatize():
    data = request.json
    doc = nlp(data["text"])
    lemmas = [token.lemma_ for token in doc]
    return jsonify(lemmas)

@app.route("/stem", methods=["POST"])
def stem():
    data = request.json
    words = data["text"].split()
    stems = [stemmer.stem(word) for word in words]
    return jsonify(stems)

@app.route("/classify_stopwords", methods=["POST"])
def classify_stopwords():
    data = request.json
    words = data["text"].split()
    classification = {word: (word in stop_words) for word in words}
    return jsonify(classification)

if __name__ == "__main__":
    app.run()
