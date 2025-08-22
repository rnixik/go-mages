var sceneConfigPreloader = {
    key: 'Preloader',
    preload: function() {
        this.cameras.main.backgroundColor = Phaser.Display.Color.HexStringToColor("#16181a");
        this.add.sprite(0, 0, 'menu_bg').setOrigin(0,0);

        const preloadBar = this.add.sprite(145, this.cameras.main.height - 165, 'preloaderBar').setOrigin(0,0);

        this.load.on('progress', function (value) {
            preloadBar.setCrop(0, 0, preloadBar.width * value, preloadBar.height);
        });


        this.load.image('menu_panel', 'assets/menu_panel.png');
        this.load.image('button_play', 'assets/button_play.png');
        this.load.image('spinner', 'assets/spinner.png');
        this.load.image('black', 'assets/black.png');
        this.load.image('defeat', 'assets/defeat.png');
        this.load.image('victory', 'assets/victory.png');

        this.load.image('bg', 'assets/temples.png');
        this.load.image('platform', 'assets/platform140.png');
        this.load.spritesheet('mage', 'assets/mage.png', { frameWidth: 130, frameHeight: 110 });
        this.load.spritesheet('sparks', 'assets/sparks256.png', { frameWidth: 256, frameHeight: 256 });
        this.load.image('hb_bar', 'assets/blood_red_bar.png');
        this.load.image('hb_health1', 'assets/blood_red_bar_health.png');
        this.load.image('hb_health2', 'assets/blood_red_bar_health.png');
        this.load.image('ice', 'assets/comet.png');
        this.load.image('icon_fireball', 'assets/icons/fire256.png');
        this.load.image('icon_lightning', 'assets/icons/lightning256.png');
        this.load.image('icon_earth', 'assets/icons/earth256.png');
        this.load.image('icon_ice', 'assets/icons/ice256.png');
        this.load.image('icon_frame', 'assets/icons/frame256.png');
        this.load.image('icon_frame_shield', 'assets/icons/shield256.png');
        this.load.spritesheet('lightning', 'assets/lightning.png', { frameWidth: 196, frameHeight: 534 });
        this.load.spritesheet('fireball', 'assets/fireball.png', { frameWidth: 64, frameHeight: 64 });
        this.load.spritesheet('rocks', 'assets/rocks_animation_120.png', { frameWidth: 120, frameHeight: 172 });
        this.load.spritesheet('shield_animation', 'assets/shield_animation.png', { frameWidth: 256, frameHeight: 256 });
        this.load.audio('shield_reflected', 'assets/shield_reflected.mp3');

    },
    create: function() {
        this.anims.create({
            key: 'lightning',
            frames: 'lightning',
            frameRate: 12,
            repeat: -1
        });
        this.anims.create({
            key: 'fireball',
            frames: 'fireball',
            frameRate: 12,
            repeat: -1
        });
        this.anims.create({
            key: 'sparks',
            frames: 'sparks',
            frameRate: 12,
            repeat: -1
        });
        this.anims.create({
            key: 'rocks',
            frames: 'rocks',
            frameRate: 12,
            repeat: -1
        });
        this.anims.create({
            key: 'shield_animation',
            frames: 'shield_animation',
            frameRate: 12,
            repeat: -1
        });

        this.scene.switch('MainMenu');
    }
};
