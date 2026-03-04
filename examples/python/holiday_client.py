/""
Holiday API Indonesia - Python Example

This example demonstrates how to use the Holiday API with Python
using the requests library.
"""

import requests
from typing import List, Dict, Optional
from dataclasses import dataclass
from datetime import datetime


@dataclass
class Holiday:
    """Represents an Indonesian holiday"""
    id: int
    date: str
    name: str
    type: str
    description: str
    created_at: Optional[str] = None
    updated_at: Optional[str] = None

    @classmethod
    def from_dict(cls, data: Dict) -> 'Holiday':
        return cls(
            id=data.get('id'),
            date=data.get('date'),
            name=data.get('name'),
            type=data.get('type'),
            description=data.get('description'),
            created_at=data.get('created_at'),
            updated_at=data.get('updated_at')
        )


class HolidayAPIClient:
    """Client for the Holiday API Indonesia"""
    
    def __init__(self, base_url: str = "http://localhost:8080/api/v1", api_key: str = None):
        self.base_url = base_url
        self.api_key = api_key
        self.session = requests.Session()
        
        if api_key:
            self.session.headers.update({'X-API-Key': api_key})
    
    def _request(self, method: str, endpoint: str, **kwargs) -> Dict:
        """Make a request to the API"""
        url = f"{self.base_url}{endpoint}"
        headers = kwargs.pop('headers', {})
        headers['Accept'] = 'application/json'
        
        response = self.session.request(method, url, headers=headers, **kwargs)
        response.raise_for_status()
        return response.json()
    
    def get_holidays(self, year: int = None, month: int = None, 
                     holiday_type: str = None) -> List[Holiday]:
        """Get holidays with optional filters"""
        params = {}
        if year:
            params['year'] = year
        if month:
            params['month'] = month
        if holiday_type:
            params['type'] = holiday_type
        
        result = self._request('GET', '/holidays', params=params)
        if result.get('success'):
            return [Holiday.from_dict(h) for h in result.get('data', [])]
        return []
    
    def get_holidays_by_year(self, year: int) -> List[Holiday]:
        """Get all holidays for a specific year"""
        result = self._request('GET', f'/holidays/year/{year}')
        if result.get('success'):
            return [Holiday.from_dict(h) for h in result.get('data', [])]
        return []
    
    def get_holidays_by_month(self, year: int, month: int) -> List[Holiday]:
        """Get holidays for a specific year and month"""
        result = self._request('GET', f'/holidays/month/{year}/{month}')
        if result.get('success'):
            return [Holiday.from_dict(h) for h in result.get('data', [])]
        return []
    
    def get_today_holiday(self) -> Optional[Holiday]:
        """Get today's holiday if any"""
        result = self._request('GET', '/holidays/today')
        if result.get('success') and result.get('data'):
            return Holiday.from_dict(result['data'][0])
        return None
    
    def get_upcoming_holidays(self, limit: int = 10) -> List[Holiday]:
        """Get upcoming holidays"""
        result = self._request('GET', f'/holidays/upcoming?limit={limit}')
        if result.get('success'):
            return [Holiday.from_dict(h) for h in result.get('data', [])]
        return []
    
    def get_holidays_this_year(self) -> List[Holiday]:
        """Get holidays for the current year"""
        result = self._request('GET', '/holidays/this-year')
        if result.get('success'):
            return [Holiday.from_dict(h) for h in result.get('data', [])]
        return []
    
    def get_holidays_this_month(self) -> List[Holiday]:
        """Get holidays for the current month"""
        result = self._request('GET', '/holidays/this-month')
        if result.get('success'):
            return [Holiday.from_dict(h) for h in result.get('data', [])]
        return []
    
    def health_check(self) -> Dict:
        """Check API health"""
        url = self.base_url.replace('/api/v1', '') + '/health'
        response = self.session.get(url)
        response.raise_for_status()
        return response.json()


def main():
    """Example usage of the Holiday API client"""
    client = HolidayAPIClient()
    
    # Example 1: Get holidays for 2024
    print("=== Holidays in 2024 ===")
    holidays = client.get_holidays_by_year(2024)
    for holiday in holidays:
        print(f"{holiday.date}: {holiday.name} ({holiday.type})")
    
    # Example 2: Get holidays for January 2024
    print("\n=== Holidays in January 2024 ===")
    holidays = client.get_holidays_by_month(2024, 1)
    for holiday in holidays:
        print(f"{holiday.date}: {holiday.name}")
    
    # Example 3: Get today's holiday
    print("\n=== Today's Holiday ===")
    today = client.get_today_holiday()
    if today:
        print(f"Today is a holiday: {today.name}")
    else:
        print("Today is not a holiday")
    
    # Example 4: Get upcoming holidays
    print("\n=== Upcoming Holidays (next 5) ===")
    upcoming = client.get_upcoming_holidays(5)
    for holiday in upcoming:
        print(f"{holiday.date}: {holiday.name}")
    
    # Example 5: Get holidays with filters
    print("\n=== National Holidays in 2024 ===")
    national = client.get_holidays(year=2024, holiday_type='national')
    for holiday in national:
        print(f"{holiday.date}: {holiday.name}")
    
    # Example 6: Health check
    print("\n=== Health Check ===")
    health = client.health_check()
    print(f"API Status: {health.get('status')}")


if __name__ == '__main__':
    main()
