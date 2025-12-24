<div align="center">
<img src="./frontend/src/assets/logo_small.svg" width="256" alt="" />
  <h1> Medovukha 🍯</h1>
  <p> Open-source web interface for Docker </p>
</div>

## 🔧 Install (from source)
 
### 🐳 Docker
```bash
git clone https://github.com/Szent7/medovukha.git
cd medovukha
#choose standalone or compose
docker build --tag medovukha .
docker run -d --restart=always -p 10015:10015 -v /var/run/docker.sock:/var/run/docker.sock --name medovukha medovukha:latest
#OR
docker compose up -d