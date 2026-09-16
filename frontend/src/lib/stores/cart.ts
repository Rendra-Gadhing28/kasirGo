import { writable, derived } from 'svelte/store';
import type { CartItem, Customer, Product } from '$lib/types';

interface CartState {
  items: CartItem[];
  customer: Customer | null;
  discountType: 'none' | 'percentage' | 'fixed';
  discountValue: number;
  taxEnabled: boolean;
  notes: string;
}

const initialState: CartState = {
  items: [],
  customer: null,
  discountType: 'none',
  discountValue: 0,
  taxEnabled: true,
  notes: ''
};

function createCartStore() {
  const { subscribe, set, update } = writable<CartState>(initialState);

  return {
    subscribe,

    addItem: (product: Product, qty: number = 1) => {
      update((state) => {
        const existingIndex = state.items.findIndex((item) => item.product.id === product.id);
        if (existingIndex > -1) {
          const newItems = [...state.items];
          const newQty = newItems[existingIndex].qty + qty;
          if (newQty <= product.stock) {
            newItems[existingIndex].qty = newQty;
          }
          return { ...state, items: newItems };
        } else {
          if (product.stock >= qty) {
            return {
              ...state,
              items: [
                ...state.items,
                { product, qty, discount_type: 'none', discount_value: 0 }
              ]
            };
          }
        }
        return state;
      });
    },

    updateQty: (productId: number, qty: number) => {
      update((state) => {
        if (qty <= 0) {
          return {
            ...state,
            items: state.items.filter((item) => item.product.id !== productId)
          };
        }
        return {
          ...state,
          items: state.items.map((item) => {
            if (item.product.id === productId) {
              const safeQty = Math.min(qty, item.product.stock);
              return { ...item, qty: safeQty };
            }
            return item;
          })
        };
      });
    },

    removeItem: (productId: number) => {
      update((state) => ({
        ...state,
        items: state.items.filter((item) => item.product.id !== productId)
      }));
    },

    setItemDiscount: (productId: number, type: 'none' | 'percentage' | 'fixed', value: number) => {
      update((state) => ({
        ...state,
        items: state.items.map((item) => {
          if (item.product.id === productId) {
            return { ...item, discount_type: type, discount_value: value };
          }
          return item;
        })
      }));
    },

    setGlobalDiscount: (type: 'none' | 'percentage' | 'fixed', value: number) => {
      update((state) => ({ ...state, discountType: type, discountValue: value }));
    },

    toggleTax: () => {
      update((state) => ({ ...state, taxEnabled: !state.taxEnabled }));
    },

    setCustomer: (customer: Customer | null) => {
      update((state) => ({ ...state, customer }));
    },

    setNotes: (notes: string) => {
      update((state) => ({ ...state, notes }));
    },

    clearCart: () => {
      set({ ...initialState });
    }
  };
}

export const cart = createCartStore();

// Derived calculations
export const cartTotals = derived(cart, ($cart) => {
  // 1. Subtotal of items after item-level discounts
  let subtotal = 0;
  let totalItemsCount = 0;

  $cart.items.forEach((item) => {
    totalItemsCount += item.qty;
    const gross = item.product.sell_price * item.qty;
    let disc = 0;
    if (item.discount_type === 'percentage' && item.discount_value > 0) {
      disc = gross * (item.discount_value / 100);
    } else if (item.discount_type === 'fixed' && item.discount_value > 0) {
      disc = item.discount_value;
    }
    subtotal += Math.max(0, gross - disc);
  });

  // 2. Global transaction discount
  let globalDiscountAmount = 0;
  if ($cart.discountType === 'percentage' && $cart.discountValue > 0) {
    globalDiscountAmount = subtotal * ($cart.discountValue / 100);
  } else if ($cart.discountType === 'fixed' && $cart.discountValue > 0) {
    globalDiscountAmount = $cart.discountValue;
  }
  globalDiscountAmount = Math.min(globalDiscountAmount, subtotal);

  const discountedSubtotal = subtotal - globalDiscountAmount;

  // 3. Tax 11% (PPN)
  let taxAmount = 0;
  if ($cart.taxEnabled) {
    taxAmount = Math.round(discountedSubtotal * 0.11);
  }

  // 4. Final Total
  const total = discountedSubtotal + taxAmount;

  return {
    subtotal,
    globalDiscountAmount,
    taxAmount,
    total,
    itemCount: totalItemsCount
  };
});
