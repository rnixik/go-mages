var sceneConfigBoot = {
    key: 'boot',
    preload: function() {
        this.load.image('menu_bg', 'assets/menu_bg.png');
        this.load.image('preloaderBar', 'assets/loader_bar.png');
    },
    create: function() {
        this.scene.switch('Preloader');
    }
};
