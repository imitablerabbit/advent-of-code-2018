#!/usr/bin/env escript

main([Filename]) ->
    try
        {ok, Polymer} = read_file(Filename),
        Reduced = reduce_polymer(Polymer),
        %io:format("Starting Polymer: ~s~n", [Polymer]),
        io:format("Reduced Polymer: ~s~n", [Reduced]),
        file:write_file("output.data", Reduced)
    catch
        E:R ->
            io:format("Error: ~p:~p~n~p~n", [E, R, erlang:get_stacktrace()]),
            usage()
    end;
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
    %reduce_polymer(Polymer, {[], 0}).
    parallel_reduce(Polymer, 4).

parallel_reduce(Polymer, WorkerNum) ->
    Size = round(length(Polymer) / WorkerNum),
    Self = self(),
    [spawn(fun() ->
                   Chunk = lists:sublist(Polymer, (N * Size) + 1, Size),
                   Reduced = reduce_polymer(Chunk, {[], 0}),
                   Self ! {N, Reduced}
           end) || N <- lists:seq(0, WorkerNum-1)],
    Combined = [receive
        {N, Reduced} ->
             Reduced
     end || N <- lists:seq(0, WorkerNum-1)],
    reduce_polymer(lists:flatten(Combined), {[], 0}).

-spec reduce_polymer(Polymer :: string(), Acc :: {Reduced :: string(), Changes :: pos_integer()}) -> string().
reduce_polymer([], {Reduced, Changes}) when Changes =:= 0 ->
    lists:reverse(Reduced);
reduce_polymer([U1], {Reduced, Changes}) when Changes =:= 0 ->
    lists:reverse([U1 | Reduced]);
reduce_polymer([], {Reduced, Changes}) when Changes > 0 ->
    io:format("Reductions: ~p, Length: ~p~n", [Changes, length(Reduced)]),
    reduce_polymer(lists:reverse(Reduced), {[], 0});
reduce_polymer([U1], {Reduced, Changes}) when Changes > 0 ->
    io:format("Reductions: ~p, Length: ~p~n", [Changes, length(Reduced)]),
    reduce_polymer(lists:reverse([U1 | Reduced]), {[], 0});
reduce_polymer([U1, U2 | T], {Reduced, Changes}) ->
    case should_reduce(U1, U2) of
        true ->
            reduce_polymer(T, {Reduced, Changes + 1});
        false ->
            reduce_polymer([U2 | T], {[U1 | Reduced], Changes})
    end.

%% They can be reduced is the strings are equal when comparing case insensitive and
%% when the units are not directly equal.
should_reduce(U, U) ->
    true;
should_reduce(U1, U2) ->
    LU1 = string:lowercase([U1]),
    LU2 = string:lowercase([U2]),
    LU1 == LU2.

