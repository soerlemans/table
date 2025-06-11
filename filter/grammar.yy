// Values
%token IDENTIFIER NUMBER STRING

// Logical:
%token TRUE FALSE

%token '!'
%token OR AND

%token COMMA
%token PIPE // Acts like a terminating symbol.
%token DOT
%token DOLLAR_SIGN

// Keywords:
%token  WHEN MUT
%token  SORT NUMERIC_SORT
%token  HEAD TAIL

// Write directives.
%token CSV MD JSON HTML

// Comparison:
%token LE LTE EQ NE GE GTE

%start program
%%
accessor_name    : DOT IDENTIFIER
								 | DOT STRING
								 ;

accessor_positional : DOLLAR_SIGN IDENTIFIER
								    | DOLLAR_SIGN NUMBER
								    ;

column           : accessor_name
								 | accessor_positional
                 ;

parameter        : rvalue
                 ;

parameter_list   : // empty
                 | parameter_list COMMA parameter
                 | parameter
                 ;

rvalue           : IDENTIFIER
                 | NUMBER
                 | STRING
								 | column
                 ;

expr             : rvalue LE rvalue
                 | rvalue LTE rvalue
                 | rvalue EQ rvalue
                 | rvalue NE rvalue
                 | rvalue GE rvalue
                 | rvalue GTE rvalue
                 ;

keyword          : WHEN expr
                 | MUT expr
								 | SORT rvalue
								 | NUMERIC_SORT rvalue
								 | HEAD rvalue
								 | TAIL rvalue
                 ;

format_csv        : CSV
                 | CSV parameter_list
                 ;

format_md         : MD
                 | MD parameter_list
                 ;

format_pretty     : PRETTY
                 | PRETTY parameter_list
                 ;

format_json       : JSON
                 | JSON parameter_list
                 ;

format_html       : HTML
                 | HTML parameter_list
                 ;

format_out       : format_csv
                 | format_md
                 | format_pretty
                 | format_json
                 | format_html
                 ;

write_directive  : GT rvalue
                 ;


// Items:
item             : keyword
                 | expr
                 ;

item_list        : // empty
                 | item_list PIPE item
                 | item
                 | item_list PIPE format_out // A format_out should always be at the end (if it is there).
                 ;

program          : item_list
								 | format_out
                 ;

%%
