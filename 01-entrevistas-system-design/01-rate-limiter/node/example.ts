import { TokenBucket } from './token-bucket';

// Simulamos una base de datos de configuración de clientes
interface ClientConfig {
  id: string;
  tier: 'free' | 'pro' | 'enterprise';
  bucket: TokenBucket;
}

// Configuración de Tiers
const TIER_LIMITS = {
  free: { capacity: 2, refillRate: 0.5 },      // 2 burst, 1 req cada 2 seg
  pro: { capacity: 5, refillRate: 4 },        // 5 burst, 5 req/seg
  enterprise: { capacity: 50, refillRate: 20 } // 50 burst, 20 req/seg
};

// Inicializamos clientes
const clients: ClientConfig[] = [
  {
    id: 'client-free-01',
    tier: 'free',
    bucket: new TokenBucket(TIER_LIMITS.free.capacity, TIER_LIMITS.free.refillRate)
  },
  {
    id: 'client-pro-01',
    tier: 'pro',
    bucket: new TokenBucket(TIER_LIMITS.pro.capacity, TIER_LIMITS.pro.refillRate)
  }
];

function simulateTraffic() {
  console.log('🚀 Iniciando simulación de tráfico...\n');

  const totalSteps = 10;
  let step = 0;

  const interval = setInterval(() => {
    step++;
    const now = new Date().toISOString().split('T')[1].replace('Z', '');

    console.log(`--- T=${step}s [${now}] ---`);

    clients.forEach(client => {
      // Intentamos hacer 3 peticiones simultáneas por cliente en cada tick
      let successful = 0;
      const attempts = 5;

      for (let i = 0; i < attempts; i++) {
        if (client.bucket.allowRequest()) {
          successful++;
        }
      }

      const statusIcon = successful === attempts ? '✅' : (successful > 0 ? '⚠️' : '❌');

      console.log(
        `${statusIcon} ${client.id} (${client.tier}): ` +
        `Procesados ${successful}/${attempts} | ` +
        `Tokens restantes: ${client.bucket.getTokens().toFixed(2)}`
      );
    });

    console.log(''); // Espacio vacío

    if (step >= totalSteps) {
      clearInterval(interval);
      console.log('🏁 Simulación terminada.');
    }
  }, 1000); // Un "tick" de simulación cada segundo real
}

simulateTraffic();