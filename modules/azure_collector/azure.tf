terraform {
  required_providers {
    bindplane = {
      source = "BlueMedoraPublic/bindplane"
      version = "0.2.4"
    }
  }
}

resource "random_uuid" "admin_password" {
}

resource "random_id" "suffix" {
  byte_length = 5
}

data "template_file" "user_data" {
  template = file("${path.module}/../../scripts/${var.userdata_script}")

  vars = {
    api_key        = var.bindplane_secret_key
    name = "${var.collector_name_prefix}-${random_id.suffix.hex}"
  }
}

resource "azurerm_public_ip" "collector" {
  name                      = "${var.collector_name_prefix}-${random_id.suffix.hex}-pubint"
  location                     = var.location
  resource_group_name          = var.resource_group
  public_ip_address_allocation = "static"
}

resource "azurerm_network_interface" "collector" {
  name                      = "${var.collector_name_prefix}-${random_id.suffix.hex}-int"
  location            = var.location
  resource_group_name = var.resource_group

  ip_configuration {
    name                      = "${var.collector_name_prefix}-${random_id.suffix.hex}-cfg"
    subnet_id                     = var.subnet_id
    private_ip_address_allocation = "dynamic"
    public_ip_address_id          = azurerm_public_ip.collector.id
  }
}

resource "azurerm_virtual_machine" "collector" {
  name                      = "${var.collector_name_prefix}-${random_id.suffix.hex}"
  location              = var.location
  resource_group_name   = var.resource_group
  network_interface_ids = [azurerm_network_interface.collector.id]
  vm_size               = var.compute_instance_size

  delete_os_disk_on_termination = var.delete_os_disk_on_termination

  storage_image_reference {
    publisher = var.compute_image["publisher"]
    offer     = var.compute_image["offer"]
    sku       = var.compute_image["sku"]
    version   = var.compute_image["version"]
  }

  storage_os_disk {
    name              = "${var.collector_name_prefix}-${random_id.suffix.hex}-os"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = var.disk_type
  }

  os_profile {
    computer_name  = "${var.collector_name_prefix}-${random_id.suffix.hex}-osprof"
    admin_username = var.admin_username
    admin_password = random_uuid.admin_password.result
    custom_data    = data.template_file.user_data.rendered
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "bindplane_collector" "collector" {
  name                      = "${var.collector_name_prefix}-${random_id.suffix.hex}"
  depends_on = [azurerm_virtual_machine.collector]
}
