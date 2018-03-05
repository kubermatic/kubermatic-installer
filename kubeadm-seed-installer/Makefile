include ../rules.make

install: cloud.conf config.sh
	true # TODO invoke install.sh (it isn't idempotent)

config.sh: ../infrastructure/osk-cluster/variables.makotemplate ../infrastructure/osk-cluster/terraform.tfstate

# TODO re-enable this check
# cloud.conf:
# 	if [ -z "$${OS_AUTH_URL}" -o -z "$${OS_USERNAME}" ]; then \
# 	    echo "OpenStack access credentials not found in environment. Please set this up first."; \
# 	    exit 1; \
# 	fi
