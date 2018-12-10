#!/usr/bin/env escript

main([Filename]) ->
        {ok, Polymer} = read_file(Filename),
        Reduced = reduce_polymer(Polymer),
        io:format("Units: ~p~n", [length(Reduced)]),
        Units = [begin
                     NewPolymer = lists:filter(
                                    fun(Elem) ->
                                        not string:equal([Elem], [Char], true)
                                    end, Reduced),
                     NewReduced = reduce_polymer(NewPolymer),
                     {length(NewReduced), [Char]}
                 end || Char <- "abcdefghijklmnopqrstuvwxyz"],
        io:format("Improved Polymer Units: ~p~n", [lists:min(Units)]),
        file:write_file("output.data", Reduced);
main([]) ->
    usage().

usage() ->
    io:format("Usage: polymer.escript 'filename'~n").

-spec read_file(Filename :: string()) -> {ok, string()}.
read_file(Filename) ->
    {ok, Bin} = file:read_file(Filename),
    Polymer = binary_to_list(Bin),
    {ok, Polymer}.

-spec reduce_polymer(Polymer :: string()) -> string().
reduce_polymer(Polymer) ->
    reduce_polymer(Polymer, []).
    
reduce_polymer([], Reduced) ->
    lists:reverse(Reduced);
reduce_polymer([U1], Reduced) ->
    lists:reverse([U1 | Reduced]);
reduce_polymer([U1, U2 | T], Reduced) ->
    case should_reduce(U1, U2) of
        true ->
            case Reduced of 
                [] ->
                    reduce_polymer(T, []);
                [RH1 | NewReduced] ->    
                    reduce_polymer([RH1 | T], NewReduced)
            end;
        false ->
            reduce_polymer([U2 | T], [U1 | Reduced])
    end.

%% They can be reduced is the strings are equal when comparing case insensitive and
%% when the units are not directly equal.
should_reduce(U, U) ->
    false;
should_reduce(U1, U2) ->
    string:equal([U1], [U2], true).
