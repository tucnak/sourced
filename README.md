# Sourced
>Sourced is a preprocessor for Source (game engine) scripts, written in Go.

Writing complex scripts for games like Team Fortress 2 usually might be really frustrating, especially when it comes down to sophisticated binds. Source engine also does not provide any sort of type checks, so scripts are hard to write, debug and maintain.

Sourced lets you write some good looking binds and produces valid output. This
```d
// Masking uber.
bind R {
	say_team "~~~ Ready ~~~"
	voicemenu 1 1
}

// Faking uber.
bind F {
	say_team "~~~ Faked ~~~"
	voicemenu 1 7
}

// Spawn forwarding.
alias switch_random {
	join_class random
	bind BACKSPACE switch_medic
}

alias switch_medic {
	join_class medic
	bind BACKSPACE switch_random
	say_team "~~~ Switched spawns ~~~"
}

bind BACKSPACE switch_random

bind K {
	slot2
	kill
}

// Altering between heal guns.
bind "[" {
	load_itempreset 0
	say_team "~~~ Switched to uber ~~~"
}

bind "]" {
	load_itempreset 1
	say_team "~~~ Took critz ~~~"
}

bind "\\" {
	load_itempreset 2
	say_team "~~~ Took vaccinator ~~~"
}

echo "~~~ Arrr! Dat healz! ~~~"
```

becomes this:
```d
alias bind_R "say_team "~~~ Ready ~~~";voicemenu 1 1";bind R bind_R
alias bind_F "say_team "~~~ Faked ~~~";voicemenu 1 7";bind F bind_F
alias switch_random "join_class random;alias bind_BACKSPACE switch_medic;bind BACKSPACE bind_BACKSPACE"
alias switch_medic "join_class medic;alias bind_BACKSPACE switch_random;bind BACKSPACE bind_BACKSPACE;say_team "~~~ Switched spawns ~~~""
alias bind_BACKSPACE switch_random;bind BACKSPACE bind_BACKSPACE
alias bind_K "slot2;kill";bind K bind_K
alias bind_LBRACKET "load_itempreset 0;say_team "~~~ Switched to uber ~~~"";bind "[" bind_LBRACKET
alias bind_RBRACKER "load_itempreset 1;say_team "~~~ Took critz ~~~"";bind "]" bind_RBRACKER
alias bind_BACKSLASH "load_itempreset 2;say_team "~~~ Took vaccinator ~~~"";bind "\" bind_BACKSLASH
echo "~~~ Arrr! Dat healz! ~~~"
```

Sourced would also do a type check for every command call and emit a type error, for instance if `alias` command takes more than two parameters.
