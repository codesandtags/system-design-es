defmodule ElixirPoc.MixProject do
  use Mix.Project

  def project do
    [
      app: :elixir_poc,
      version: "0.1.0",
      elixir: "~> 1.19",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger]
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:jason, "~> 1.4"}, # El parser de JSON más rápido y estándar en Elixir
      {:flow, "~> 1.2"}   # La librería que nos permite usar todos los núcleos (Map-Reduce)
    ]
  end
end
