# TypeScript cheat sheet

## Modelado de Datos y validaciones básicas

El 90% de muchos problemas manejan diccionarios. (ej. Pais -> Balance).

**Para objetos definidios (Entidades)**: Es mejor usar `inteface`. Es limpio y legible.

```typescript
interface Account {
  country: string;
  balance: number;
  currency?: string; // optional
}
```

**Para Mapas/Diccionarios (Agregaciones)**: No uses `any`, mejor usa `Record<key, value>`. Una mala practica es usar `{ [key: string]: any }` porque pierdes tipado.

```typescript
const balances: Record<string, number> = {};

// safe assignments
balances["USA"] = 100;
balances["CAN"] = 200;
```

**Parseo robusto**: Maneja `unknown` y `type guards`. Evita a toda costa `any`, si parseas algo desconocido, usa `unknown` y valida su tipo antes de usarlo. Si escribes una funcion pequeña que valida tipos, ganas puntos instantaneos de `Design Approach`.

```typescript
// Example of user-defined type guard in TypeScript.
interface Account {
  country: string;
  balance: number;
  currency?: string;
}

function isAccount(obj: unknown): obj is Account {
  if (typeof obj !== "object" || obj === null) {
    return false;
  }

  const account = obj as Account;
  return (
    typeof account.country === "string" &&
    typeof account.balance === "number" &&
    (account.currency === undefined || typeof account.currency === "string")
  );
}

// Usage
const data: unknown = JSON.parse(someString);
if (isAccount(data)) {
  console.log("Valid Account:", data.country, data.balance);
}
```

**Array Magic: reduce es tu mejor amigo**: Para problemas de "mover dinero" o "calcular totales", reduce es más elegante que un loop for.

```typescript
// Example of using reduce to calculate total balance
interface Account {
  country: string;
  balance: number;
  currency?: string;
}

const accounts: Account[] = [
  { country: "USA", balance: 100 },
  { country: "CAN", balance: 200 },
  { country: "MEX", balance: 150 },
  { country: "USA", balance: 50 },
  { country: "USA", balance: 30 },
];

const balanceMap = accounts.reduce<Record<string, number>>((acc, account) => {
  // Accumulate balance per country. If country not in acc, initialize to 0.
  acc[account.country] = (acc[account.country] || 0) + account.balance;
  return acc;
}, {});
```

**Manejo de Nulos y `undefined` strictness**: TypeScript generara errores si tratas de acceder a propiedades que pueden ser `undefined` sin validarlas primero. Para evitar errores en tiempo de ejecucion, siempre valida o usa el operador de encadenamiento opcional `?.` o el operador de asercion no nula `!`.

- Usa ! SOLO si estás 100% seguro (ej. en tests). map.get('US')!.
- Usa ?. y ?? (Nullish Coalescing) en lógica de negocio.

```typescript
const usBalance = balances["USA"] ?? 0; // Si no existe, usa 0
const total = usBalance + 100; // safe to use
```

**Utility Types**: TypeScript ofrece varios tipos utilitarios que pueden ser muy útiles para manipular tipos existentes. Algunos de los más comunes incluyen:

- `Partial<T>`: Hace que todas las propiedades de T sean opcionales.
- `Required<T>`: Hace que todas las propiedades de T sean obligatorias.
- `Readonly<T>`: Hace que todas las propiedades de T sean de solo lectura.
- `Pick<T, K>`: Crea un nuevo tipo seleccionando un subconjunto de propiedades de T.
- `Omit<T, K>`: Crea un nuevo tipo omitiendo un subconjunto de propiedades de T.
- `Record<K, T>`: Crea un tipo de objeto con claves de tipo K y valores de tipo T.

Algo importante a resltar resaltar es que para el segundo argumento se puede usar una unión de tipos.

```typescript
interface Account {
  country: string;
  balance: number;
  currency?: string;
}

type AccountSummary = Pick<Account, "country" | "balance">;
const summary: AccountSummary = {
  country: "USA",
  balance: 100,
};

// Function that updates an account with partial updates.
function updateAccount(account: Account, updates: Partial<Account>): Account {
  return { ...account, ...updates };
}
```
