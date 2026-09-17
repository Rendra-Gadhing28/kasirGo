export interface User {
  id: number;
  outlet_id: number;
  name: string;
  email: string;
  role: 'owner' | 'admin' | 'kasir';
  is_active: boolean;
  last_login?: string;
  created_at: string;
}

export interface Outlet {
  id: number;
  name: string;
  address: string;
  phone: string;
  logo?: string;
}

export interface Category {
  id: number;
  outlet_id: number;
  name: string;
  description: string;
  created_at: string;
}

export interface Product {
  id: number;
  outlet_id: number;
  category_id: number;
  category?: Category;
  name: string;
  sku: string;
  barcode: string;
  buy_price: number;
  sell_price: number;
  stock: number;
  unit: string;
  image_url: string;
  is_active: boolean;
  created_at: string;
}

export interface Customer {
  id: number;
  outlet_id: number;
  member_code?: string;
  name: string;
  phone: string;
  email: string;
  address: string;
  points: number;
  created_at: string;
}

export interface Supplier {
  id: number;
  outlet_id: number;
  name: string;
  phone: string;
  email: string;
  address: string;
  contact_person: string;
  created_at: string;
}

export interface CartItem {
  product: Product;
  qty: number;
  discount_type: 'none' | 'percentage' | 'fixed';
  discount_value: number;
}

export interface TransactionItem {
  id: number;
  transaction_id: number;
  product_id: number;
  product_name: string;
  product_sku: string;
  qty: number;
  unit: string;
  price: number;
  buy_price: number;
  discount_type: string;
  discount_value: number;
  discount_amount: number;
  subtotal: number;
}

export interface Transaction {
  id: number;
  outlet_id: number;
  outlet?: Outlet;
  invoice_no: string;
  customer_id?: number;
  customer?: Customer;
  user_id: number;
  user?: User;
  subtotal: number;
  discount_type: string;
  discount_value: number;
  discount_amount: number;
  tax_rate: number;
  tax_amount: number;
  total_amount: number;
  payment_method: 'cash' | 'qris' | 'transfer';
  payment_status: 'pending' | 'paid' | 'cancelled' | 'refunded';
  paid_amount: number;
  change_amount: number;
  notes: string;
  midtrans_order_id?: string;
  midtrans_qr_code?: string;
  midtrans_transaction_status?: string;
  items: TransactionItem[];
  created_at: string;
}

export interface APIResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  meta?: {
    current_page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
  errors?: any;
}
