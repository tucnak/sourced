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
```d
with ALT {
    bind F {
        say "hello"
    }
}
```

This obviously looks really promisisng when the whole thing comes down to tangled scripts, like vaccinator workaround:
```d
// Vaccinator binds, simplyfing switches:
//
// [Z] against bullets
// [X] against explosives
// [C] against fire
//
// You should not use [R] to change vaccine,
// since it would break script functionality.
alias next_vaccine reload

// Defaults.
alias against_bullets ""

alias against_explosives {
    next_vaccine
}

alias against_fire {
    next_vaccine
    next_vaccine
}

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
    // ...
}

bind C {
    // ...
}

echo "~~~ Arrr! Dat healz! ~~~"
```

Sourced would also do a type check for every command call and emit a type error, for instance if `alias` command takes more than two parameters.
