from flask import Flask
import os

app = Flask(__name__)

@app.route("/")
def hello_world():
    value = os.getenv()
    return "<p>Hello, World!</p>"
