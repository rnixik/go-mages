const prepareDuration = 300;

const GameScene = function() {
    this.socket = null;
    this.joinedData = {};

    this.myPlayerIndex = 0;
    this.myClientId = 0;
    this.player1 = null;
    this.player2 = null ;

    this.spellButtons = [];
    this.castedSpells = [];
    this.game = null;
    this.sendGameCommand = function () {};
};

GameScene.prototype = {
    create: function(game, data) {
        const self = this;

        this.myClientId = data.myClientId;
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
            1);
        this.player2 = new Player(
            this.game,
            game.cameras.main.width - platformSidePadding - platformWidth / 2,
            game.cameras.main.height - platformHeight - platformWBottomPadding,
            2);

        this.player1.draw();
        this.player2.draw();
        this.player1.stateDefault();
        this.player2.stateDefault();

        this.spellButtons = this.game.add.group();
        this.addSpellIcon('fireball', 15, 60, 'icon_fireball');
        this.addSpellIcon('rocks', 94, 60, 'icon_earth');
        this.addSpellIcon('lightning', 15, 140, 'icon_lightning');
        this.addSpellIcon('comet', 94, 140, 'icon_ice');

        this.addSpellProtectionIcon('protect_fireball', game.cameras.main.width - 140, 50, 'icon_fireball');
        this.addSpellProtectionIcon('protect_rocks', game.cameras.main.width - 75, 50, 'icon_earth');
        this.addSpellProtectionIcon('protect_lightning', game.cameras.main.width - 140, 125, 'icon_lightning');
        this.addSpellProtectionIcon('protect_comet', game.cameras.main.width - 75, 125, 'icon_ice');
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
        var spellIcon = this.game.add.sprite(x, y, spellIconResourceId).setOrigin(0, 0).setDisplaySize(64, 64);
        this.game.add.sprite(x, y, 'icon_frame').setOrigin(0, 0).setDisplaySize(64, 64);
        spellIcon.spellId = spellId;
        spellIcon.setInteractive();
        var self = this;
        spellIcon.on('pointerdown', function () {
            self.onIconClick(spellIcon);
        });
    },
    addSpellProtectionIcon: function(spellId, x, y, spellIconResourceId) {
        this.game.add.sprite(x + 20, y + 20, spellIconResourceId).setOrigin(0, 0).setDisplaySize(20, 20);
        let shieldIcon = this.game.add.sprite(x, y, 'icon_frame_shield').setOrigin(0, 0).setDisplaySize(64, 64);
        shieldIcon.spellId = spellId;
        shieldIcon.setInteractive();
        let self = this;
        shieldIcon.on('pointerdown', function () {
            self.onIconClick(shieldIcon);
        });
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
        let spell;
        switch (data.spellId) {
            case 'fireball':
                spell = new FireballSpell();
                break;
            case 'lightning':
                spell = new LightningSpell();
                castingDuration = 500;
                break;
            case 'comet':
                spell = new CometSpell();
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
        originPlayer.statePreparation(data.spellId);
        let _this = this;
        setTimeout(function() {
            spell.cast(_this.game, originPlayer, targetPlayer);
            originPlayer.stateCasting(castingDuration);
        }, prepareDuration);

    },
    onPlayersUpdate: function(data) {
        var pl1 = this.player1;
        var pl2 = this.player2;
        if (this.myPlayerIndex === 2) {
            pl1 = this.player2;
            pl2 = this.player1;
        }
        pl1.setHealthBar(data.hp1 / 100);
        pl2.setHealthBar(data.hp2 / 100);
    },
    onEndGame: function(data) {
        this.add.sprite(0, 0, 'black').alpha = 0.9;
        if (data.winner === this.myPlayerIndex) {
            var image = this.add.sprite(0, 100, 'victory');
        } else {
            var image = this.add.sprite(0, 100, 'defeat');
        }

        var _this = this;
        setTimeout(function() {
            image.inputEnabled = true;
            image.events.onInputDown.add(function() {
                _this.game.state.start('MainMenu');
            });
        }, 1000);

    },
    onIncomingGameEvent: function (name, data) {
        switch (name) {
            case 'CastEvent':
                this.onCast(data);
                break;
        }
    }
};


function Player(game, posX, posY, playerIndex) {
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
}

Player.prototype.draw = function() {
    this.sprite = this.game.add.sprite(this.posX, this.posY, this.resourceId, 1);
    this.sprite.setOrigin(0.5, 1);
    if (this.index === 2) {
        this.sprite.scaleX = -1;
    }

    this.drawHealthBar();
    this.setHealthBar(1);
};

Player.prototype.getSpellPos = function() {
    var spellPosition = {x: this.posX, y: this.posY - 25};
    return spellPosition;
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

Player.prototype.stateCasting = function(duration) {
    var _this = this;
    setTimeout(function(){
        _this._changeAnimation(1, duration);
    }, 10);
};

Player.prototype.statePreparation = function(spellId) {
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
