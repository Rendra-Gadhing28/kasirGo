export interface BarcodeLookupResult {
  barcode: string;
  name: string;
  brand?: string;
  sku?: string;
  imageUrl?: string;
  unit?: string;
  categoryHint?: string;
  found: boolean;
  raw?: any;
}

function detectSuggestedUnit(quantity: string, categories: string): string {
  const q = (quantity || '').toLowerCase();
  const c = (categories || '').toLowerCase();

  if (q.includes('ml') || q.includes('liter') || q.includes(' l') || c.includes('beverage') || c.includes('water') || c.includes('drink')) {
    return 'botol';
  }
  if (c.includes('canned') || c.includes('can') || q.includes('kaleng')) {
    return 'kaleng';
  }
  if (q.includes('kg') || q.includes('gram') || q.includes(' gr') || q.includes(' g') || c.includes('noodle') || c.includes('snack') || c.includes('biscuit')) {
    return 'bungkus';
  }
  if (q.includes('sachet') || c.includes('sachet')) {
    return 'sachet';
  }
  return 'pcs';
}

function detectCategoryHint(categories: string, tags: string[] = []): string {
  const allText = `${categories} ${tags.join(' ')}`.toLowerCase();

  if (allText.includes('beverage') || allText.includes('drink') || allText.includes('water') || allText.includes('tea') || allText.includes('coffee') || allText.includes('juice') || allText.includes('soda') || allText.includes('milk') || allText.includes('minuman')) {
    return 'Makanan & Minuman';
  }
  if (allText.includes('noodle') || allText.includes('mie') || allText.includes('rice') || allText.includes('beras') || allText.includes('oil') || allText.includes('minyak') || allText.includes('sugar') || allText.includes('gula') || allText.includes('sembako')) {
    return 'Sembako';
  }
  if (allText.includes('snack') || allText.includes('biscuit') || allText.includes('chip') || allText.includes('wafer') || allText.includes('candy') || allText.includes('cokelat') || allText.includes('chocolate') || allText.includes('camilan') || allText.includes('keripik')) {
    return 'Snack & Cemilan';
  }
  if (allText.includes('shampoo') || allText.includes('soap') || allText.includes('sabun') || allText.includes('bath') || allText.includes('toothpaste') || allText.includes('dental') || allText.includes('skin') || allText.includes('mandi') || allText.includes('perawatan')) {
    return 'Perawatan & Mandi';
  }
  if (allText.includes('detergent') || allText.includes('cleaner') || allText.includes('tissue') || allText.includes('tisu') || allText.includes('pembersih') || allText.includes('rumah')) {
    return 'Kebutuhan Rumah';
  }

  return 'Makanan & Minuman';
}

function generateSmartSKU(brand: string, name: string, barcode: string): string {
  if (brand && brand.length >= 2) {
    const cleanBrand = brand.replace(/[^A-Za-z0-9]/g, '').slice(0, 3).toUpperCase();
    const lastDigits = barcode.slice(-4);
    return `${cleanBrand}-${lastDigits}`;
  }
  return `SKU-${barcode.slice(-6)}`;
}

async function fetchFromOFF(url: string): Promise<any | null> {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 3500);

  try {
    const res = await fetch(url, {
      signal: controller.signal
    });
    if (res.ok) {
      const json = await res.json();
      if (json.status === 1 && json.product) {
        return json.product;
      }
    }
  } catch (err: any) {
    // Aborted or network error
  } finally {
    clearTimeout(timeoutId);
  }
  return null;
}

export async function lookupBarcodeInfo(barcode: string): Promise<BarcodeLookupResult | null> {
  const cleanCode = barcode.trim().replace(/[^0-9A-Za-z]/g, '');
  if (!cleanCode) return null;

  const fields = 'product_name,product_name_id,product_name_en,generic_name,brands,brand_owner,quantity,image_front_url,image_url,categories,categories_tags';

  // 1. Prioritize Indonesian Open Food Facts instance, fallback to World instance
  const urlsToTry = [
    `https://id.openfoodfacts.org/api/v2/product/${cleanCode}?fields=${fields}`,
    `https://world.openfoodfacts.org/api/v2/product/${cleanCode}?fields=${fields}`
  ];

  // If 12 digits, also prepare 13-digit EAN query with leading zero
  if (cleanCode.length === 12) {
    urlsToTry.push(`https://id.openfoodfacts.org/api/v2/product/0${cleanCode}?fields=${fields}`);
    urlsToTry.push(`https://world.openfoodfacts.org/api/v2/product/0${cleanCode}?fields=${fields}`);
  }

  for (const url of urlsToTry) {
    const p = await fetchFromOFF(url);
    if (p) {
      const brand = (p.brands || p.brand_owner || '').split(',')[0].trim();
      
      // Indonesian name prioritised, then generic / general name
      const primaryName = (p.product_name_id || p.product_name || p.product_name_en || p.generic_name || '').trim();
      const quantity = (p.quantity || '').trim();
      const categories = p.categories || '';
      const tags = p.categories_tags || [];

      // Construct clean, descriptive supermarket product title
      let fullName = primaryName;
      if (brand && !fullName.toLowerCase().includes(brand.toLowerCase())) {
        fullName = `${brand} ${fullName}`.trim();
      }
      if (quantity && !fullName.toLowerCase().includes(quantity.toLowerCase())) {
        fullName = `${fullName} ${quantity}`.trim();
      }

      const suggestedUnit = detectSuggestedUnit(quantity, categories);
      const categoryHint = detectCategoryHint(categories, tags);
      const smartSKU = generateSmartSKU(brand, fullName, cleanCode);

      return {
        barcode: cleanCode,
        name: fullName || `Produk Barcode ${cleanCode}`,
        brand: brand,
        sku: smartSKU,
        imageUrl: p.image_front_url || p.image_url || '',
        unit: suggestedUnit,
        categoryHint: categoryHint,
        found: true,
        raw: p
      };
    }
  }

  // 2. Secondary Fast Fallback: UPC Item DB (1.5s timeout)
  try {
    const upcController = new AbortController();
    const upcTimeout = setTimeout(() => upcController.abort(), 2000);

    const upcRes = await fetch(`https://api.upcitemdb.com/prod/trial/lookup?upc=${cleanCode}`, {
      signal: upcController.signal
    });
    clearTimeout(upcTimeout);

    if (upcRes.ok) {
      const upcData = await upcRes.json();
      if (upcData.items && upcData.items.length > 0) {
        const item = upcData.items[0];
        const brand = (item.brand || '').trim();
        const smartSKU = generateSmartSKU(brand, item.title, cleanCode);

        return {
          barcode: cleanCode,
          name: item.title || `Produk Barcode ${cleanCode}`,
          brand: brand,
          sku: smartSKU,
          imageUrl: (item.images && item.images.length > 0) ? item.images[0] : '',
          unit: 'pcs',
          categoryHint: 'Makanan & Minuman',
          found: true,
          raw: item
        };
      }
    }
  } catch (e) {
    // ignore
  }

  // Not found in databases: Return clean fallback for manual input
  return {
    barcode: cleanCode,
    name: '',
    brand: '',
    sku: `SKU-${cleanCode.slice(-6)}`,
    unit: 'pcs',
    categoryHint: '',
    found: false
  };
}
