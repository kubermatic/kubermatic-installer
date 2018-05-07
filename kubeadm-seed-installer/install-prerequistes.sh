#!/bin/bash
which cfssl || sudo curl -o /usr/local/bin/cfssl https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
which cfssljson || sudo curl -o /usr/local/bin/cfssljson https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
[[ -x /usr/local/bin/cfssl ]] || sudo chmod +x /usr/local/bin/cfssl
[[ -x /usr/local/bin/cfssljson ]] || sudo chmod +x /usr/local/bin/cfssljson
