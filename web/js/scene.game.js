const GameScene = function() {
    this.socket = null;
    this.joinedData = {};

    this.myPlayerIndex = 0;
    this.player1 = null;
    this.player2 = null ;

    this.spellButtons = [];
    this.castedSpells = [];
    this.game = null;
};

GameScene.prototype = {
    create: function(game) {
        this.game = game;
        this.game.add.sprite(0, 0, 'bg').setOrigin(0, 0);

        var platforms = this.game.add.group();

        var platformTextureImage = this.game.textures.get('platform').getSourceImage();
        var platformWidth = platformTextureImage.width;
        var platformHeight = platformTextureImage.height;
        var platformWBottomPadding = 30;
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


        this.spellButtons = this.game.add.group();
        this.addSpellIcon('fireball', 15, 60, 'icon_fireball', 'icon_frame_red');
        this.addSpellIcon('lightning', 15, 140, 'icon_lightning', 'icon_frame_blue');
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
    addSpellIcon: function(spellId, x, y, spellIconResourceId, spellFrameResourceId) {
        var spellIcon = this.game.add.sprite(x, y, spellIconResourceId).setOrigin(0, 0).setDisplaySize(64, 64);
        this.game.add.sprite(x, y, spellFrameResourceId).setOrigin(0, 0).setDisplaySize(64, 64);
        spellIcon.spellId = spellId;
        spellIcon.setInteractive();
        var self = this;
        spellIcon.on('pointerdown', function () {
            self.onIconClick(spellIcon);
        });

    },
    onIconClick: function(icon) {
        if (icon.spellId) {
            // this.socket.emit('action', {'action': 'cast', 'spell': icon.spellId});
        }
        this.player1.stateCasting(400);
    },
    onCast: function(data) {
        var spell = false;
        if (data.spell === 'lightning') {
            spell = new LightningSpell(this);
        } else if (data.spell === 'fireball') {
            spell = new FireballSpell(this);
        }

        var targetPlayer = this.player1;
        var sourcePlayer = this.player2;
        if (data.sourcePlayer === this.myPlayerIndex) {
            sourcePlayer = this.player1;
            targetPlayer = this.player2;
        }

        if (spell) {
            spell.cast(sourcePlayer, targetPlayer);
            sourcePlayer.stateCasting(spell.playerAnimateDuration);
            this.castedSpells.push(spell);
            //console.log(data.spell, data.sourcePlayer, this.myPlayerIndex);
        }

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
    this.stateDefault();
    var _this = this;
    setTimeout(function(){
        _this._changeAnimation(1, duration);
    }, 10);
};

Player.prototype._changeAnimation = function(frame, duration) {
    this.sprite.setFrame(frame);
    console.log(this.sprite.anims.currentFrame);
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
    create: function () {
        gameScene.create(this);
    },
    update: function () {
        gameScene.update(this);
    },
};
