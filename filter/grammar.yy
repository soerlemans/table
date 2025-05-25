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
%token  WHEN MUT OUT MD JSON

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

keyword_ out     : OUT
                 | OUT STRING
                 | OUT STRING COMMA parameter_list
                 ;

keyword          : WHEN expr
                 | MUT expr
                 | keyword_out
                 | MD parameter_list
                 | JSON parameter_list
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

// Items:
item             : keyword
                 | expr
                 ;

item_list        : // empty
                 | item_list PIPE item
                 | item
                 ;

program          : item_list
                 ;

%%
