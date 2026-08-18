import os
from flask import Flask
from dotenv import load_dotenv

from auth import auth_bp, init_oauth
from routes import routes_bp

load_dotenv()

def create_app():
    app = Flask(__name__)
    app.secret_key = os.getenv("FLASK_SECRET_KEY")

    init_oauth(app)

    app.register_blueprint(auth_bp)
    app.register_blueprint(routes_bp)

    return app

app = create_app()

if __name__ == "__main__":
    # This is the port specified in okta 
    app.run(debug=True, port=8080)


