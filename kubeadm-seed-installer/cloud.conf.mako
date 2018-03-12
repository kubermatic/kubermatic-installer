<%!
import os
import sys
%>
<%
if not('OS_USERNAME' in os.environ) or not('OS_AUTH_URL' in os.environ):
    sys.exit("OpenStack access credentials not found in environment. Please set this up first.")
%>
[Global]
auth-url = "${os.environ['OS_AUTH_URL']}"
username = "${os.environ['OS_USERNAME']}"
password = "${os.environ['OS_PASSWORD']}"
domain-name = "Default"
tenant-name = "${os.environ['OS_TENANT_NAME']}"
region = "${os.environ['OS_REGION_NAME']}"

[BlockStorage]
trust-device-path = false
bs-version = "v2"
ignore-volume-az=true 
