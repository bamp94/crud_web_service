create or replace function test.user_get(_id integer)
    returns json as
$BODY$
declare
    _ret json;
begin
    if _id = 0 then
        select array_to_json(array(
                select row_to_json(r)
                from (
                         select u.id, u.name, u.email
                         from test.users u
                     ) r
            ))
        into _ret;
    else
        select row_to_json(r)
        into _ret
        from (
                 select u.id, u.name, u.email
                 from test.users u
                 where id = _id
             ) r;
    end if;

    return _ret;

exception
    when others then
        return json_build_object('error', SQLERRM);
end
$BODY$
    language plpgsql volatile
                     cost 100;


create or replace function test.user_ins(_params json)
    returns json as
$BODY$
declare
    _newid integer;
begin
    _newid = 0;

    insert into test.users (name, email)
    select name, email
    from json_populate_record(null::test.users, _params)
    returning id into _newid;

    return json_build_object('id', _newid);

exception
    when others then
        return json_build_object('error', SQLERRM);
end
$BODY$
    language plpgsql volatile
                     cost 100;


create or replace function test.user_upd(_id integer, _params json)
    returns json as
$BODY$
begin
    update test.users
    set name  = _params ->> 'name',
        email = _params ->> 'email'
    where id = _id;

    return json_build_object('id', _id);

exception
    when others then
        return json_build_object('error', SQLERRM);
end
$BODY$
    language plpgsql volatile
                     cost 100;


create or replace function test.user_del(_id integer)
    returns json as
$BODY$
begin
    delete from test.users where id = _id;

    return json_build_object('id', _id);

exception
    when others then
        raise notice 'Illegal operation: %', SQLERRM;

        return json_build_object('error', SQLERRM);
end
$BODY$
    language plpgsql volatile
                     cost 100;