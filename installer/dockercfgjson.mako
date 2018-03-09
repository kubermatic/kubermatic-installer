<%! import os %>
<%! from base64 import b64encode %>

<%include file="variables.mako" />

% if os.path.isfile("variables_override.mako"):
<%include file="variables_override.mako" />
% endif

<%
var._dockercfgjson = b64encode('{"auths":{"https://index.docker.io/v1/":{"username":"%s","password":"%s","email":"%s","auth":"%s"}}}' % (var.docker_username, var.docker_password, var.docker_email, b64encode(var.docker_username + ":" + var.docker_password)))
%>
