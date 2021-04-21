FOO = b"""# HELP foo foo
# TYPE foo counter
foo{service_instance_id="850cba52-ac6a-430c-8f00-16945785c0c2"} 1000
"""
BAR = b"""
# HELP bar bar
# TYPE bar counter
bar{service_instance_id="850cba52-ac6a-430c-8f00-16945785c0c2"} 1000
"""


def app(environ, start_response):
    start_response("200 OK", (("Content-Length", str(len(FOO) + len(BAR))),))
    yield FOO

    if app.CLOSE_SOCKET:
        # Force close the socket via gunicorn
        sock = environ["gunicorn.socket"]
        sock.close()
    else:
        yield BAR

    # Invert state on every request
    app.CLOSE_SOCKET = not app.CLOSE_SOCKET


app.CLOSE_SOCKET = False
