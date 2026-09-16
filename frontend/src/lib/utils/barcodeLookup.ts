export interface BarcodeLookupResult {
  barcode: string;
  name: string;
  brand?: string;
  imageUrl?: string;
  unit?: string;
  raw?: any;
}

export async function lookupBarcodeInfo(barcode: string): Promise<BarcodeLookupResult | null> {
  const cleanCode = barcode.trim();
  if (!cleanCode) return null;

  try {
    // 1. Primary Source: Open Food Facts API (Extensive Indonesian FMCG & Food database)
    const offRes = await fetch(`https://world.openfoodfacts.org/api/v0/product/${cleanCode}.json`, {
      headers: {
        'User-Agent': 'KasirPro-POS - Web - Version 1.0'
      }
    });

    if (offRes.ok) {
      const data = await offRes.json();
      if (data.status === 1 && data.product) {
        const p = data.product;
        const brand = p.brands || p.brand_owner || '';
        let productName = p.product_name || p.product_name_id || p.product_name_en || p.generic_name || '';

        // Prepend brand if not already included in title
        let fullName = productName;
        if (brand && !productName.toLowerCase().includes(brand.toLowerCase())) {
          fullName = `${brand} ${productName}`.trim();
        }

        // Detect suggested unit
        let suggestedUnit = 'pcs';
        const qtyStr = (p.quantity || '').toLowerCase();
        if (qtyStr.includes('ml') || qtyStr.includes('l') || qtyStr.includes('liter')) {
          suggestedUnit = 'botol';
        } else if (qtyStr.includes('kg') || qtyStr.includes('gram') || qtyStr.includes('g')) {
          suggestedUnit = 'bungkus';
        }

        return {
          barcode: cleanCode,
          name: fullName || `Produk Barcode ${cleanCode}`,
          brand: brand,
          imageUrl: p.image_front_url || p.image_url || '',
          unit: suggestedUnit,
          raw: p
        };
      }
    }
  } catch (e) {
    console.warn('Open Food Facts lookup failed:', e);
  }

  // 2. Fallback: UPC Item DB (General retail/goods database)
  try {
    const upcRes = await fetch(`https://api.upcitemdb.com/prod/trial/lookup?upc=${cleanCode}`);
    if (upcRes.ok) {
      const upcData = await upcRes.json();
      if (upcData.items && upcData.items.length > 0) {
        const item = upcData.items[0];
        return {
          barcode: cleanCode,
          name: item.title || `Produk Barcode ${cleanCode}`,
          brand: item.brand || '',
          imageUrl: (item.images && item.images.length > 0) ? item.images[0] : '',
          unit: 'pcs',
          raw: item
        };
      }
    }
  } catch (e) {
    console.warn('UPC Item DB lookup failed:', e);
  }

  // Return minimal info if barcode is recognized but details not in database
  return {
    barcode: cleanCode,
    name: `Produk Baru (${cleanCode})`,
    brand: '',
    unit: 'pcs'
  };
}
