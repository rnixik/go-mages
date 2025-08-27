const GameScene = function() {
    this.socket = null;
    this.joinedData = {};

    this.myPlayerIndex = 1;
    this.myClientId = 0;
    this.myNickname = 'Unknown';
    this.opponentNickname = 'Unknown';
    this.player1 = null;
    this.player2 = null ;

    this.isDesktop = false;
    this.spellPressed = false;
    this.shieldPressed = false;

    this.spellButtons = [];
    this.castedSpells = [];
    this.game = null;
    this.sendGameCommand = function () {};

    this.attackCooldownMs = 2000;
    this.shieldCooldownMs = 900;
    this.maxHp = 1000;
};

GameScene.prototype = {
    create: function(game, data) {
        const self = this;

        this.isDesktop = game.scene.systems.game.device.os.desktop;
        this.myClientId = data.myClientId;
        this.myNickname = data.myNickname;
        this.opponentNickname = data.opponentNickname;
        console.log("Game started. My client id: " + this.myClientId + ", my nickname: " + this.myNickname + ", opponent nickname: " + this.opponentNickname);
        this.sendGameCommand = data.sendGameCommand;
        data.setOnIncomingGameEventCallback(function (name, data) {
            self.onIncomingGameEvent(name, data);
        });

        this.game = game;
        this.game.add.sprite(0, 0, 'bg').setOrigin(0, 0);

        var platforms = this.game.add.group();

        var platformTextureImage = this.game.textures.get('platform').getSourceImage();
        var platformWidth = platformTextureImage.width;
        var platformHeight = platformTextureImage.height;
        var platformWBottomPadding = 39;
        var platformSidePadding = 80;

        var platform1 = platforms.create(
            platformSidePadding,
            this.game.cameras.main.height - platformHeight - platformWBottomPadding,
            'platform'
        ).setOrigin(0, 0);

        var platform2 = platforms.create(
            game.cameras.main.width - platformSidePadding - platformWidth / 2,
            game.cameras.main.height - platformHeight - platformWBottomPadding,
            'platform'
        ).setOrigin(0, 0);

        platform2.setOrigin(0.5, 0);
        platform2.scale.x *= -1;


        this.player1 = new Player(
            this.game,
            platformSidePadding + platformWidth / 2,
            game.cameras.main.height - platformHeight - platformWBottomPadding,
            1,
            this.myNickname
        );
        this.player2 = new Player(
            this.game,
            game.cameras.main.width - platformSidePadding - platformWidth / 2,
            game.cameras.main.height - platformHeight - platformWBottomPadding,
            2,
            this.opponentNickname
        );

        this.player1.draw();
        this.player2.draw();
        this.player1.stateDefault();
        this.player2.stateDefault();

        this.spellButtons = this.game.add.group();
        this.addSpellIcon('fireball', 136, 120, 'icon_fireball');
        this.addSpellIcon('rocks', 76, 180, 'icon_earth');
        this.addSpellIcon('lightning', 76, 60, 'icon_lightning');
        this.addSpellIcon('comet', 16, 120, 'icon_ice');


        this.addSpellProtectionIcon('protect_fireball', game.cameras.main.width - 16 - 62, 120, 'icon_fireball');
        this.addSpellProtectionIcon('protect_rocks', game.cameras.main.width - 76 - 62, 180, 'icon_earth');
        this.addSpellProtectionIcon('protect_lightning', game.cameras.main.width - 76 - 62, 60, 'icon_lightning');
        this.addSpellProtectionIcon('protect_comet', game.cameras.main.width - 136 - 62, 120, 'icon_ice');
        //
        // this.onConnected();
    },
    update: function() {
        // for (var s in this.castedSpells) {
        //     if (this.castedSpells[s].isActive) {
        //         this.castedSpells[s].onUpdate();
        //     }
        // }
    },
    onConnected: function() {

        // this.myPlayerIndex = this.joinedData.playerIndex;

        // var _this = this;
        // this.socket.on('state', function(data) {
        //     if (data.state === 'update') {
        //         _this.onPlayersUpdate(data);
        //     } else if (data.state === 'cast') {
        //         _this.onCast(data);
        //     } else if (data.state === 'endGame') {
        //         _this.onEndGame(data);
        //     }
        // });
        //
        // this.socket.on('disconnect', function() {
        //     alert('Connection lost');
        // });
    },
    addSpellIcon: function(spellId, x, y, spellIconResourceId) {
        var spellIcon = this.game.add.sprite(x, y, spellIconResourceId).setOrigin(0, 0).setDisplaySize(60, 60);
        this.game.add.sprite(x, y, 'icon_frame').setOrigin(0, 0).setDisplaySize(60, 60);

        let frame = null;

        spellIcon.spellId = spellId;
        spellIcon.setInteractive();
        var self = this;
        let onSpellClicked = function () {
            if (self.spellPressed) {
                return;
            }
            self.spellPressed = true;
            self.onIconClick(spellIcon);
            frame = self.game.add.sprite(x+4, y+4, 'black_frame').setOrigin(0, 0).setDisplaySize(52, 52);
            spellIcon.setPosition(x+4, y+4).setDisplaySize(52, 52);

            self.cooldownButton(spellIcon, self.attackCooldownMs);

            self.game.time.delayedCall(self.attackCooldownMs, () => {
                spellIcon.setPosition(x, y).setDisplaySize(60, 60);
                frame && frame.destroy(true, true);
                self.spellPressed = false;
            });
        };

        spellIcon.on('pointerdown', onSpellClicked);

        if (this.isDesktop) {
            let keyFrameIndex = 0;
            let keyName = '';
            switch (spellId) {
                case 'lightning':
                    keyFrameIndex = 3; // W
                    keyName = 'W';
                    break;
                case 'comet':
                    keyFrameIndex = 0; // A
                    keyName = 'A';
                    break;
                case 'rocks':
                    keyFrameIndex = 1; // S
                    keyName = 'S';
                    break;
                case 'fireball':
                    keyFrameIndex = 2; // D
                    keyName = 'D';
                    break;
            }
            this.game.add.sprite(x + 60 - 24 - 2, y + 60 - 24 - 2, 'keys', keyFrameIndex).setOrigin(0, 0).setDisplaySize(24, 24);

            this.game.scene.scene.input.keyboard.on("keydown-" + keyName, function (event) {
                onSpellClicked();
            });
        }
    },
    addSpellProtectionIcon: function(spellId, x, y, spellIconResourceId) {
        let shieldIcon = this.game.add.sprite(x, y, 'icon_frame_shield').setOrigin(0, 0).setDisplaySize(62, 62);
        let spellIcon = this.game.add.sprite(x + 16, y + 16, spellIconResourceId).setOrigin(0, 0).setDisplaySize(30, 30);
        this.game.add.sprite(x + 16, y + 16, 'icon_frame').setOrigin(0, 0).setDisplaySize(30, 30);
        shieldIcon.spellId = spellId;
        shieldIcon.setInteractive();
        let self = this;
        let onShieldClicked = function () {
            if (self.shieldPressed) {
                return;
            }
            self.shieldPressed = true;
            self.onIconClick(shieldIcon);
            spellIcon.setPosition(x + 16 + 4, y + 16 + 4).setDisplaySize(30 - 4, 30 - 4);
            self.cooldownButton(shieldIcon, self.shieldCooldownMs);
            
            self.game.time.delayedCall(self.shieldCooldownMs, () => {
                spellIcon.setPosition(x+16, y+16).setDisplaySize(30, 30);
                self.shieldPressed = false;
            });
        };

        shieldIcon.on('pointerdown', onShieldClicked);

        if (this.isDesktop) {
            let keyFrameIndex = 0;
            let keyName = '';
            switch (spellId) {
                case 'protect_lightning':
                    keyFrameIndex = 7;
                    keyName = 'I';
                    break;
                case 'protect_comet':
                    keyFrameIndex = 4;
                    keyName = 'J';
                    break;
                case 'protect_rocks':
                    keyFrameIndex = 5;
                    keyName = 'K';
                    break;
                case 'protect_fireball':
                    keyFrameIndex = 6;
                    keyName = 'L';
                    break;
            }
            this.game.add.sprite(x + 60 - 24 - 2, y + 60 - 24 - 2, 'keys', keyFrameIndex).setOrigin(0, 0).setDisplaySize(24, 24);

            this.game.scene.scene.input.keyboard.on("keydown-" + keyName, function (event) {
                onShieldClicked();
            });
        }
    },
    onIconClick: function(icon) {
        if (icon.spellId) {
            this.sendGameCommand('Cast', {'spellId': icon.spellId})
        }
    },
    onCast: function(data) {
        let targetPlayer = this.player1;
        let originPlayer = this.player2;
        if (this.myClientId === data.OriginPlayerId) {
            originPlayer = this.player1;
            targetPlayer = this.player2;
        }

        let castingDuration = 400;
        let prepareDuration = 300;
        let spell;
        switch (data.spellId) {
            case 'fireball':
                spell = new FireballSpell();
                prepareDuration = 200;
                break;
            case 'lightning':
                spell = new LightningSpell();
                castingDuration = 500;
                break;
            case 'comet':
                spell = new CometSpell();
                prepareDuration = 200;
                break;
            case 'rocks':
                spell = new RocksSpell();
                break;
            case 'protect_fireball':
            case 'protect_lightning':
            case 'protect_rocks':
            case 'protect_comet':
                spell = new ShieldSpell();
                spell.cast(this.game, originPlayer, data.spellId);
                originPlayer.stateShield(300);
                return;
            default:
                throw Error("Unknown spell: " + data.spellId);
        }

        originPlayer.stateDefault();
        originPlayer.statePreparation(data.spellId, 400);
        let _this = this;
        setTimeout(function() {
            spell.cast(_this.game, originPlayer, targetPlayer);
            originPlayer.stateCasting(castingDuration);
        }, prepareDuration);

    },
    onDamage: function(data) {
        if (data.shieldWorked) {
            this.game.sound.play('shield_reflected');
        }
        let targetPlayer = this.player2;
        if (this.myClientId === data.targetPlayerId) {
            targetPlayer = this.player1;
        }
        targetPlayer.setHealthBar(data.targetPlayerHp / this.maxHp);
    },
    onEndGame: function(data) {
        this.game.add.sprite(0, 0, 'black').setOrigin(0, 0).alpha = 0.9;
        let image;
        if (this.myClientId === data.winnerPlayerId) {
            image = this.game.add.sprite(0, 100, 'victory').setOrigin(0, 0);
        } else {
            image = this.game.add.sprite(0, 100, 'defeat').setOrigin(0, 0);
        }

        var _this = this;
        setTimeout(function() {
            image.setInteractive();
            image.on('pointerdown', function () {
                image.destroy(true, true);
                _this.game.scene.switch('MainMenu');
            });

        }, 1000);

    },
    onIncomingGameEvent: function (name, data) {
        switch (name) {
            case 'CastEvent':
                this.onCast(data);
                break;
            case 'DamageEvent':
                this.onDamage(data);
                break;
            case 'EndGameEvent':
                this.onEndGame(data);
                break;
        }
    },
    cooldownButton: function (sprite, ms) {
        sprite.disableInteractive();
        sprite.setTint(0x808080);
        this.game.time.delayedCall(ms, () => {
            sprite.setInteractive();
            sprite.clearTint();
        });
    }
};


function Player(game, posX, posY, playerIndex, nickname) {
    this.game = game;
    this.posX = posX;
    this.posY = posY;
    this.index = playerIndex; //1 or 2
    this.sprite = false;

    this.hbSprite = false;
    this.hbSpriteBaseWidth = 0;
    this.hbBasePosX = 0;

    this.resourceId = 'mage';
    this.changeAnimationTimeoutId = null;

    this.nickname = nickname;
}

Player.prototype.draw = function() {
    this.sprite = this.game.add.sprite(this.posX, this.posY, this.resourceId, 1);
    this.sprite.setOrigin(0.5, 1);
    if (this.index === 2) {
        this.sprite.scaleX = -1;
    }

    this.drawHealthBar();
    this.setHealthBar(1);
    this.drawNickname();
};

Player.prototype.getSpellPos = function() {
    return {x: this.posX, y: this.posY - 25};
};

Player.prototype.drawHealthBar = function() {
    var healthBar = this.game.add.group();
    var resourceId = 'hb_health1';
    this.hbBasePosX = 10;
    if (this.index === 2) {
        resourceId = 'hb_health2';
        this.hbBasePosX = this.game.cameras.main.width - 10 - this.game.textures.get(resourceId).getSourceImage().width;
    }
    this.hbSprite = healthBar.create(this.hbBasePosX, 10, resourceId).setOrigin(0, 0);
    healthBar.create(this.hbBasePosX, 10, 'hb_bar').setOrigin(0, 0);
    this.hbSpriteBaseWidth = this.hbSprite.width;

};

Player.prototype.drawNickname = function() {
    this.game.make.text({
        x: this.posX,
        y: this.posY + 60,
        text: this.nickname,
        style: {
            fontFamily: 'Arial',
            color: '#ffffff',
            align: 'center',
        },
        origin: {x: 0.5, y: 0.5},
        add: true
    });
};

Player.prototype.stateCasting = function(duration) {
    var _this = this;
    setTimeout(function(){
        _this._changeAnimation(1, duration);
    }, 10);
};

Player.prototype.statePreparation = function(spellId, prepareDuration) {
    this.sprite.setFrame(2);
    let sprite = this.game.add.sprite(this.posX - 130, this.posY - 180, 'sparks', 1).setOrigin(0, 0);
    sprite.play('sparks');

    // set color
    switch (spellId) {
        case 'fireball':
            sprite.tint = 0xff6f00;
            break;
        case 'lightning':
            sprite.tint = 0xffd800;
            break;
        case 'comet':
            sprite.tint = 0x86bfe0;
            break;
        case 'rocks':
            sprite.tint = 0x684740;
            break;
    }

    setTimeout(function () {
        sprite.destroy(true, true);
    }, prepareDuration);
};

Player.prototype.stateShield = function(duration) {
    var _this = this;
    setTimeout(function(){
        _this._changeAnimation(3, duration);
    }, 10);
};

Player.prototype._changeAnimation = function(frame, duration) {
    this.sprite.setFrame(frame);
    clearTimeout(this.changeAnimationTimeoutId);
    var _this = this;
    this.changeAnimationTimeoutId = setTimeout(function(){
        _this.stateDefault();
    }, duration);
};

Player.prototype.stateDefault = function() {
    this.sprite.setFrame(0);
};

Player.prototype.setHealthBar = function(health) {
    health = Math.max(0, health);
    health = Math.min(1, health);
    var cropSize = parseInt(this.hbSpriteBaseWidth * health);
    this.hbSprite.setCrop(0, 0, cropSize, this.hbSprite.height);
    if (this.index === 2) {
        this.hbSprite.x = this.hbBasePosX + this.hbSpriteBaseWidth - cropSize - 5; //5 = hb's padding
        this.hbSprite.x = Math.max(this.hbSprite.x, this.hbBasePosX);
    }
};

const gameScene = new GameScene();

var sceneConfigGame = {
    key: 'Game',
    create: function (data) {
        gameScene.create(this, data);
    },
    update: function () {
        gameScene.update(this);
    },
};
