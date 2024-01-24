
# tailwindcss.exe  -i ./static/css/input.css -o static/css/output.css --watch
# tailwindcss.exe build -i ./static/css/input.css -o ./static/css/output.css --minify
templ fmt .
templ generate --watch --proxy="https://localhost:4000" --cmd "tailwindcss -i ./static/css/input.css -o static/css/output.css && go run ."
