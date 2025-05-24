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

column_list      : // empty
                 | column_list COMMA column
                 | column
                 ;

keyword          : WHEN expr
                 | MUT expr
                 | OUT STRING column_list
                 | MD column_list
                 | JSON column_list
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
