# GO Runner Online (IUT Nantes)

![Bannière](https://ibb.co/m6pzJxk)

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


