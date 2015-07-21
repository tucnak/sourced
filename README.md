# Sourced
>Sourced is a preprocessor for Source (game engine) scripts, written in Go.

Writing complex scripts for games like Team Fortress 2 usually might be really frustrating, especially when it comes down to sophisticated binds. Source engine also does not provide any sort of type checks, so scripts are hard to write, debug and maintain.

A common, yet complex routine is support of multi-key binds. For example, saying "hello" in team chat with <Alt+F> keyboard combination would require:
```
alias +alt_mod "bind F "say "hello"""
alias -alt_mod "unbind F"
bind ALT alt_mod
```

Now take a look at solution Sourced provides:
```
with ALT {
    bind F {
        say "hello"
    }
}
```

This obviously looks really promisisng when the whole thing comes down to tangled scripts, like vaccinator workaround:
```
alias next_vaccine reload

// By default, it's against bullets => no action required.
alias against_bullets ""

bind Z {
    against_bullets

    alias against_bullets ""

    // It takes one switch from bullets to explosives.
    alias against_explosives {
        next_vaccine
    }

    // It takes two switches from bullets to fire.
    alias against_fire {
        next_vaccine
        next_vaccine
    }
}

bind X {
    against_explosives

    alias against_explosives ""

    // It takes a single switch from explosives to fire.
    alias against_fire {
        next_vaccine
    }

    // It takes two switches from explosives to bullets.
    alias against_bullets {
        next_vaccine
        next_vaccine
    }
}

bind C {
    against_fire

    alias against_fire ""

    // It takes a single switch from fire to bullets.
    alias against_bullets {
        next_vaccine
    }

    // It takes two switches from fire to explosives.
    alias against_explosives {
        next_vaccine
        next_vaccine
    }
}
```

Sourced would also do a type check for every command call and emit a type error, for instance if `alias` command takes more than two parameters.
