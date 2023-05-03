# GO Runner Online (IUT Nantes)
![GO CI](https://github.com/npldevfr/go-runner-online/actions/workflows/go.yml/badge.svg)


![Bannière](https://i.ibb.co/hD08sqb/Capture-d-e-cran-2023-05-03-a-11-48-29.png)

## Description

GO Runner Online est un jeu de course en 2D qui a pour but de terminer avant les autres, la base du project fonctionne
hors connexion, le but est de faire une version en ligne.

## Installation
Build le projet avec la commande suivante :
```
go build
```

Coté serveur (lancement):
```
./course -server
```

Coté client (lancement):
```
./course -client <ipServer>:<port>
```

## TODO
- [x] Création du projet
- [x] Création du serveur
- [x] Création du client
- [ ] Création des méthodes pour communiquer entre le client et le serveur
- [ ] Synchronisation des joueurs


