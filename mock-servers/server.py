from http.server import SimpleHTTPRequestHandler, HTTPServer
import json

class CustomHandler(SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/health':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            health_response = {
                "status": "healthy",
                "timestamp": "2024-01-01T00:00:00Z"
            }
            self.wfile.write(json.dumps(health_response).encode())
        else:
            # Handle all other requests normally (serve static files)
            super().do_GET()

def run_server(port, directory):
    handler = lambda *args: CustomHandler(*args, directory=directory)
    server = HTTPServer(('', port), handler)
    print(f"Server running on port {port}")
    server.serve_forever()

if __name__ == "__main__":
    import sys
    if len(sys.argv) != 3:
        print("Usage: python custom_server.py <port> <directory>")
        sys.exit(1)

    port = int(sys.argv[1])
    directory = sys.argv[2]
    run_server(port, directory)
