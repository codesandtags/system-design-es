defmodule ElixirPocTest do
  use ExUnit.Case
  doctest ElixirPoc

  test "greets the world" do
    assert ElixirPoc.hello() == :world
  end
end
