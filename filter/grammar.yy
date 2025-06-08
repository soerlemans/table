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

write_csv        : CSV
                 | CSV STRING
                 | CSV STRING COMMA parameter_list
                 ;

write_md         : MD
                 | MD STRING
                 | MD STRING COMMA parameter_list
                 ;

write_json       : JSON
                 | JSON STRING
                 | JSON STRING COMMA parameter_list
                 ;

write_html       : HTML
                 | HTML STRING
                 | HTML STRING COMMA parameter_list
                 ;

write            : write_csv
                 | write_md
                 | write_json
                 | write_html
                 ;


// Items:
item             : keyword
                 | expr
                 ;

item_list        : // empty
                 | item_list PIPE item
                 | item
                 ;

program          : item_list
                 | item_list PIPE write // A write should always be at the end (if it is there).
								 | write
                 ;

%%
