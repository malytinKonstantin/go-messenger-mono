Для push образов в локальный registry, вам может потребоваться добавить его в список небезопасных registry в Docker. Создайте или отредактируйте файл /etc/docker/daemon.json:

{
  "insecure-registries" : ["localhost:5000"]
}

sudo systemctl restart docker
