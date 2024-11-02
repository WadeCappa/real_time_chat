

if [ ! -d "frontend/venv" ]; then
    echo "file not found"
    python -m venv frontend/venv
    source frontend/venv/bin/activate
    pip install flask
    pip install gunicorn
else
    source frontend/venv/bin/activate
fi

cd frontend
gunicorn -w 2 'hello_world:app'
