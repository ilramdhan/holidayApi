/**
 * Holiday API Indonesia - Node.js Example
 * 
 * This example demonstrates how to use the Holiday API with Node.js
 * using the native fetch API (Node.js 18+)
 */

const BASE_URL = 'http://localhost:8080/api/v1';

class HolidayAPIClient {
  constructor(baseURL, apiKey = null) {
    this.baseURL = baseURL;
    this.apiKey = apiKey;
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const headers = {
      'Accept': 'application/json',
      ...options.headers
    };

    if (this.apiKey) {
      headers['X-API-Key'] = this.apiKey;
    }

    const response = await fetch(url, {
      ...options,
      headers
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  // Get all holidays with optional filters
  async getHolidays(params = {}) {
    const queryString = new URLSearchParams(params).toString();
    const endpoint = `/holidays${queryString ? '?' + queryString : ''}`;
    return this.request(endpoint);
  }

  // Get holidays by year
  async getHolidaysByYear(year) {
    return this.request(`/holidays/year/${year}`);
  }

  // Get holidays by month
  async getHolidaysByMonth(year, month) {
    return this.request(`/holidays/month/${year}/${month}`);
  }

  // Get today's holiday
  async getTodayHoliday() {
    return this.request('/holidays/today');
  }

  // Get upcoming holidays
  async getUpcomingHolidays(limit = 10) {
    return this.request(`/holidays/upcoming?limit=${limit}`);
  }

  // Get holidays for current year
  async getHolidaysThisYear() {
    return this.request('/holidays/this-year');
  }

  // Get holidays for current month
  async getHolidaysThisMonth() {
    return this.request('/holidays/this-month');
  }

  // Health check
  async healthCheck() {
    const response = await fetch(`${this.baseURL.replace('/api/v1', '')}/health`);
    return response.json();
  }
}

// Example usage
async function main() {
  const client = new HolidayAPIClient(BASE_URL);

  try {
    // Example 1: Get holidays for 2024
    console.log('=== Holidays in 2024 ===');
    const holidays2024 = await client.getHolidaysByYear(2024);
    if (holidays2024.success) {
      holidays2024.data.forEach(h => {
        console.log(`${h.date}: ${h.name} (${h.type})`);
      });
    }

    // Example 2: Get holidays for January 2024
    console.log('\n=== Holidays in January 2024 ===');
    const janHolidays = await client.getHolidaysByMonth(2024, 1);
    if (janHolidays.success) {
      janHolidays.data.forEach(h => {
        console.log(`${h.date}: ${h.name}`);
      });
    }

    // Example 3: Get today's holiday
    console.log('\n=== Today\'s Holiday ===');
    const today = await client.getTodayHoliday();
    if (today.success && today.data.length > 0) {
      console.log('Today is a holiday:', today.data[0].name);
    } else {
      console.log('Today is not a holiday');
    }

    // Example 4: Get upcoming holidays
    console.log('\n=== Upcoming Holidays (next 5) ===');
    const upcoming = await client.getUpcomingHolidays(5);
    if (upcoming.success) {
      upcoming.data.forEach(h => {
        console.log(`${h.date}: ${h.name}`);
      });
    }

    // Example 5: Get holidays with filters
    console.log('\n=== National Holidays in 2024 ===');
    const nationalHolidays = await client.getHolidays({ 
      year: 2024, 
      type: 'national' 
    });
    if (nationalHolidays.success) {
      nationalHolidays.data.forEach(h => {
        console.log(`${h.date}: ${h.name}`);
      });
    }

    // Example 6: Health check
    console.log('\n=== Health Check ===');
    const health = await client.healthCheck();
    console.log('API Status:', health.status);

  } catch (error) {
    console.error('Error:', error.message);
  }
}

// Run examples
main();

// Export for use as a module
module.exports = { HolidayAPIClient };
