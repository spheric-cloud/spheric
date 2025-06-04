# `cloud-hypervisor` setup

1. Make sure your OS supports KVM (check via `kvm-ok`)
2. Add the user to the kvm group
   ```shell
   sudo usermod -aG kvm $USER
   newgrp kvm # Start new shell with group membership changes
   ```
