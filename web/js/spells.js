const LightningSpell = function () {};
LightningSpell.prototype = {
    cast: function(game, sourcePlayer, targetPlayer) {
        const spellPos = sourcePlayer.getSpellPos();

        let xPos = spellPos.x + 18;
        if (sourcePlayer.index === 2) {
            xPos = spellPos.x - 18;
        }

        const sprite = game.add.sprite(xPos, spellPos.y + 88, 'lightning').setOrigin(0, 0);
        sprite.rotation = -Math.PI / 2;

        if (sourcePlayer.index === 2) {
            sprite.scaleY *= -1;
        }

        sprite.play('lightning');

        setTimeout(function () {
            sprite.destroy();
        }, 500);
    }
}

const FireballSpell = function () {};
FireballSpell.prototype = {
    cast: function (game, sourcePlayer, targetPlayer) {
        const group = game.physics.add.group();
        const spellPos = sourcePlayer.getSpellPos();

        const xPos = spellPos.x;
        const sprite = group.create(xPos, spellPos.y - 72, 'fireball').setOrigin(0, 0);
        if (sourcePlayer.index === 2) {
            sprite.scaleX *= -1;
        }
        sprite.play('fireball');

        group.setVelocity(1000, 0);
        if (sourcePlayer.index === 2) {
            group.setVelocity(-1000, 0);
        }

        setTimeout(function () {
            group.destroy(true, true);
        }, 400);
    }
};
