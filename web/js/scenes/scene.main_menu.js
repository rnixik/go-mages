const MainMenu = function () {
    this.connectFailed = false;
    this.wsConnection = null;

    this.menuControls = null;
    this.menuBackground = null;
    this.loadingSpinner = null;

    const self = this;
    this.game = null;
    this.onIncomingGameEventCallback = function () {};

    this.myClientId = null;
    this.nickname = 'default';

    this.create = function(game) {
        this.game = game;
        this.menuBackground = game.add.sprite(0, 0, 'menu_bg').setOrigin(0, 0);
        this.menuControls = game.add.group();
        var menuPanel = game.add.sprite(game.cameras.main.centerX - 269, 170, 'menu_panel').setOrigin(0, 0);
        menuPanel.alpha = 0.9;
        var buttonPlay = game.add.sprite(game.cameras.main.centerX - 101, 250, 'button_play').setOrigin(0, 0).setInteractive();
        buttonPlay.on('pointerdown', function () {
            self.actionPlay();
        });
        this.menuControls.add(menuPanel);
        this.menuControls.add(buttonPlay);


        this.loadingSpinner = game.add.sprite(400, 300, 'spinner');
        this.loadingSpinner.setVisible(false);

        const savedNickname = localStorage.getItem("nickname");
        if (savedNickname) {
            this.nickname = savedNickname;
        } else {
            const defaultNickname = 'Player' + Math.floor(Math.random() * 1000);
            const inputNickname = prompt("Please enter your nickname", defaultNickname);
            if (inputNickname !== null && inputNickname.trim() !== '') {
                this.nickname = inputNickname.trim();
            } else {
                this.nickname = defaultNickname;
            }
            try {
                localStorage.setItem("nickname", this.nickname);
            } catch (e) {
                console.warn("Local storage not available, cannot save nickname");
            }
        }
        // limit nickname up to 10 chars
        this.nickname = this.nickname.substring(0, 10);
        console.log("Using nickname: " + this.nickname);
    };

    this.update = function(game) {
        this.loadingSpinner.angle += 2;
    };

    this.actionPlay = function() {
        console.log('play');
        this.menuBackground.alpha = 0.3;
        this.menuControls.setVisible(false);
        this.loadingSpinner.setVisible(true);
        this.connectToServer();
    };

    this.connectToServer = function() {
        const wsConnect = (nickname) => {
            self.wsConnection = new WebSocket(WEBSOCKET_URL);
            self.wsConnection.onopen = function () {
                self.wsConnection.send(JSON.stringify({type: 'lobby', subType: 'join', data: nickname}));
                self.wsConnection.send(JSON.stringify({type: 'lobby', subType: 'makeMatch'}));
            };
            self.wsConnection.onclose = () => {
                window.setTimeout(function () {
                    location.reload();
                }, 3000);
            };
            self.wsConnection.onmessage = function (evt) {
                const messages = evt.data.split('\n');
                for (let i = 0; i < messages.length; i++) {
                    let json;
                    try {
                        json = JSON.parse(messages[i]);
                    } catch (ex) {
                        console.warn("Json parse error", evt.data, ex);
                    }
                    if (json) {
                        self.onIncomingMessage(json, evt);
                    }
                }
            };
        };

        wsConnect(this.nickname);
    };

    this.onIncomingMessage = function (json, evt) {
        console.log('INCOMING', json);
        if (json.name === 'ClientJoinedEvent') {
            this.myClientId = json.data.yourId;
            return;
        }
        if (json.name === 'GameStartedEvent') {
            self.startGame(this.myClientId, json.data.room.members);
            return;
        }

        self.onIncomingGameEventCallback(json.name, json.data);
    };

    this.startGame = function(myClientId, players) {
        console.log('Starting game with my client id = ' + myClientId);

        let myNickname = 'Unknown';
        let opponentNickname = 'Unknown';
        for (let i = 0; i < players.length; i++) {
            if (!players[i].isPlayer) {
                continue;
            }
            if (players[i].id === myClientId) {
                myNickname = players[i].nickname;
            } else {
                opponentNickname = players[i].nickname;
            }
        }

        this.game.scene.start('Game', {
            myClientId: myClientId,
            myNickname: myNickname,
            opponentNickname: opponentNickname,
            sendGameCommand: function (type, data) {
                self.wsConnection.send(JSON.stringify({type: 'game', subType: type, data: data}));
            },
            setOnIncomingGameEventCallback: function (callback) {
                self.onIncomingGameEventCallback = callback;
            },
        });
    };
}

const mainMenu = new MainMenu();

var sceneConfigMainMenu = {
    key: 'MainMenu',
    create: function () {
        mainMenu.create(this);
    },
    update: function () {
        mainMenu.update(this);
    },
};
