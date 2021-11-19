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

        this.load.image('bg', 'assets/forest.png');
        this.load.image('platform', 'assets/platform.png');
        this.load.spritesheet('mage', 'assets/mage.png', { frameWidth: 130, frameHeight: 110 });
        this.load.image('hb_bar', 'assets/blood_red_bar.png');
        this.load.image('hb_health1', 'assets/blood_red_bar_health.png');
        this.load.image('hb_health2', 'assets/blood_red_bar_health.png');
        this.load.image('icon_fireball', 'assets/icons/fireball-red-1.png');
        this.load.image('icon_lightning', 'assets/icons/lightning-blue-1.png');
        this.load.image('icon_frame_red', 'assets/icons/frame-9-red.png');
        this.load.image('icon_frame_blue', 'assets/icons/frame-7-blue.png');
        this.load.spritesheet('lightning', 'assets/lightning.png', { frameWidth: 196, frameHeight: 534 });
        this.load.spritesheet('fireball', 'assets/fireball.png', { frameWidth: 64, frameHeight: 64 });

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

        this.scene.switch('MainMenu');
    }
};
