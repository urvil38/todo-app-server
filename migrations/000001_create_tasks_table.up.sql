CREATE TABLE IF NOT EXISTS tasks(
  id serial PRIMARY KEY,
  name VARCHAR (50) NOT NULL,
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE FUNCTION trigger_modify_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$;
COMMENT ON FUNCTION trigger_modify_updated_at IS
'FUNCTION trigger_modify_updated_at sets the value of a column named updated_at to the current timestamp.';

CREATE TRIGGER set_updated_at BEFORE INSERT OR UPDATE ON tasks
     FOR EACH ROW EXECUTE PROCEDURE trigger_modify_updated_at();
COMMENT ON TRIGGER set_updated_at ON tasks IS
'TRIGGER set_updated_at updates the value of the updated_at column to the current timestamp whenever a row is inserted or updated to the table.';