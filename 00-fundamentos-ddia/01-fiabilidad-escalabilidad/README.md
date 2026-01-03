# Capítulo 1: Fiabilidad, Escalabilidad y Mantenibilidad

> **Basado en:** _Designing Data-Intensive Applications_ de Martin Kleppmann.
> **Nivel:** Fundamental
> **Serie:** De Senior a Staff Engineer

Este documento resume los tres pilares fundamentales que distinguen a una aplicación de juguete de un sistema de producción a gran escala.

---

## 📺 Video Explicativo

[![Ver en YouTube](https://img.youtube.com/vi/TU_VIDEO_ID_AQUI/maxresdefault.jpg)](LINK_A_TU_VIDEO_AQUI)

_Haz clic en la imagen para ver la explicación completa y el dibujo en vivo._

---

## 1. 🛡️ Fiabilidad (Reliability)

**"El sistema continúa funcionando correctamente incluso cuando ocurren adversidades."**

No se trata de evitar que las cosas fallen (es imposible), sino de evitar que esos fallos rompan la experiencia del usuario.

### Conceptos Clave

- **Fallo (Fault) vs. Falla (Failure):**
  - _Fallo:_ Un componente se desvía de su especificación (ej. un disco duro muere, una excepción de red).
  - _Falla:_ El sistema entero deja de dar servicio al usuario.
  - **Meta:** Diseñar sistemas tolerantes a fallos para prevenir fallas.
- **Tipos de Fallos:**
  1.  **Hardware:** Discos, RAM, cortes de luz. (MTTF - Mean Time To Failure).
  2.  **Software:** Bugs, procesos fuera de control, cascadas de errores.
  3.  **Humanos:** La causa #1. Configuraciones erróneas, despliegues fallidos.

> **💡 Insight:** Netflix usa _Chaos Monkey_ para causar fallos intencionales y asegurar que sus sistemas automáticos de recuperación funcionen.

---

## 2. 📈 Escalabilidad (Scalability)

**"La capacidad de un sistema para manejar una carga creciente."**

La escalabilidad no es un botón mágico; es una decisión de arquitectura basada en **Parámetros de Carga**.

### ¿Cómo medimos la carga?

Depende de tu sistema:

- **Twitter / X:** Requests por segundo (Escritura vs Lectura).
- **Base de Datos:** Ratio de Lecturas/Escrituras.
- **Chat:** Usuarios concurrentes activos.

### El Caso de Estudio: Twitter (Fan-out)

El reto de Twitter no es publicar un tweet, es entregar ese tweet a los 10 millones de seguidores de un famoso en milisegundos.

```mermaid
sequenceDiagram
    participant User as Usuario (Famoso)
    participant LB as Load Balancer
    participant Service as Servicio Publicación
    participant DB as Base de Datos
    participant Cache as Redis (Timelines)

    User->>LB: Publica Tweet
    LB->>Service: Procesa Tweet
    Service->>DB: Guarda Tweet en tabla global

    rect rgb(200, 150, 255)
    note right of Service: Enfoque Fan-out (Push)
    loop Para cada Seguidor
        Service->>Cache: Inyecta Tweet en Timeline del seguidor
    end
end
```

## 3. 🔧 Mantenibilidad (Maintainability)

"La facilidad con la que diferentes ingenieros pueden trabajar en el sistema a lo largo del tiempo."

Es el pilar olvidado, pero el más costoso económicamente. Se divide en:

- **Operabilidad**: Hacerle la vida fácil al equipo de DevOps/SRE. (Buenos logs, métricas, documentación).
- **Simplicidad**: Eliminar la Complejidad Accidental. No uses Kubernetes si un VPS basta.
- **Evolubilidad**: La facilidad para hacer cambios en el futuro (Agilidad).

## El Ángulo Frontend / Fullstack

¿Por qué te debe importar esto si trabajas con React/Next.js?

| Concepto       | Aplicación en Frontend                                                                                       |
| -------------- | ------------------------------------------------------------------------------------------------------------ |
| Fiabilidad     | ¿Tu UI maneja errores con ErrorBoundaries o se pone en blanco? ¿Usas Optimistic UI para ocultar latencia?    |
| Escalabilidad  | ¿Estás bloqueando el Main Thread? ¿Usas CDN para assets estáticos? ¿Paginación vs Infinite Scroll mal hecho? |
| Mantenibilidad | "TypeScript estricto, Design Systems, Storybook. Código que tu ""yo del futuro"" entienda."                  |

## Frases / Reflexiones Clave

- The Internet was done so well that most people think of it as a natural resource like the Pacific Ocean, rather than something that was man-made. When was the last time a technology with a scale like that was so error-free? The Web, in comparison, is a joke. The Web was done by amateurs. – Alan Kay

## Recursos

- [Distributed Systems lecture series By Martin Kleppmann](https://www.youtube.com/watch?v=UEAMfLPZZhE&list=PLeKd45zvjcDFUEv_ohr_HdUFe97RItdiB)
- [Designing A Data-Intensive Future: Expert Talk • Martin Kleppmann & Jesse Anderson • GOTO 2023](https://www.youtube.com/watch?v=P-9FwZxO1zE&t=130s)
