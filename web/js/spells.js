const LightningSpell = function () {};
LightningSpell.prototype = {
    cast: function(game, sourcePlayer, targetPlayer) {
        const spellPos = targetPlayer.getSpellPos();

        let xPos = spellPos.x - 150;
        const sprite = game.add.sprite(xPos, spellPos.y - 550, 'lightning').setOrigin(0, 0);
        sprite.tint = 0xFAF550;

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


const CometSpell = function () {};
CometSpell.prototype = {
    cast: function (game, sourcePlayer, targetPlayer) {
        const group = game.physics.add.group();
        const spellPos = sourcePlayer.getSpellPos();

        const xPos = spellPos.x;
        const sprite = group.create(xPos, spellPos.y - 72, 'ice').setOrigin(0, 0);
        if (sourcePlayer.index === 2) {
            sprite.scaleX *= -1;
        }
        sprite.tint = 0xA4D8E1;

        group.setVelocity(1000, 0);
        if (sourcePlayer.index === 2) {
            group.setVelocity(-1000, 0);
        }

        setTimeout(function () {
            group.destroy(true, true);
        }, 400);
    }
};

const RocksSpell = function () {};
RocksSpell.prototype = {
    cast: function (game, sourcePlayer, targetPlayer) {
        const group = game.physics.add.group();
        const spellPos = targetPlayer.getSpellPos();

        const xPos = spellPos.x;
        const sprite = group.create(xPos - 70, spellPos.y, 'rocks').setOrigin(0, 0);
        sprite.play('rocks');

        setTimeout(function () {
            group.destroy(true, true);
        }, 400);
    }
};

const ShieldSpell = function () {};
ShieldSpell.prototype = {
    cast: function (game, sourcePlayer, shieldType) {
        const group = game.physics.add.group();
        const spellPos = sourcePlayer.getSpellPos();

        const xPos = spellPos.x;
        const sprite = group.create(xPos - 130, spellPos.y - 180, 'shield_animation').setOrigin(0, 0);

        switch (shieldType) {
            case 'protect_fireball':
                sprite.tint = 0xFF4500; // Fireball shield color
                break;
            case 'protect_lightning':
                sprite.tint = 0xFFFF00; // Lightning shield color
                break;
            case 'protect_rocks':
                sprite.tint = 0x8B4513; // Earth shield color
                break;
            case 'protect_comet':
                sprite.tint = 0xADD8E6; // Ice shield color
                break;
        }

        sprite.play('shield_animation');

        setTimeout(function () {
            group.destroy(true, true);
        }, 300);
    }
};