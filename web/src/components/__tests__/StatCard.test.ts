/**
 * StatCard 组件测试
 * 覆盖: 基础渲染、Props传递、趋势显示、样式类
 */

import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StatCard from '../Common/StatCard.vue'

// Mock Element Plus icons
const mockIcon = {
  template: '<span class="mock-icon"></span>',
}

describe('StatCard.vue', () => {
  const defaultProps = {
    title: '今日交易',
    value: '1,234',
  }

  describe('基础渲染', () => {
    it('should render title and value', () => {
      const wrapper = mount(StatCard, {
        props: defaultProps,
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.text()).toContain('今日交易')
      expect(wrapper.text()).toContain('1,234')
    })

    it('should render with default color', () => {
      const wrapper = mount(StatCard, {
        props: defaultProps,
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      const card = wrapper.find('.stat-card')
      expect(card.attributes('style')).toContain('border-left-color')
    })

    it('should apply custom color', () => {
      const wrapper = mount(StatCard, {
        props: {
          ...defaultProps,
          color: '#ff0000',
        },
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      const card = wrapper.find('.stat-card')
      expect(card.attributes('style')).toContain('#ff0000')
    })
  })

  describe('前缀和后缀', () => {
    it('should render prefix when provided', () => {
      const wrapper = mount(StatCard, {
        props: {
          ...defaultProps,
          prefix: '¥',
        },
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.prefix').exists()).toBe(true)
      expect(wrapper.find('.prefix').text()).toBe('¥')
    })

    it('should render suffix when provided', () => {
      const wrapper = mount(StatCard, {
        props: {
          ...defaultProps,
          suffix: '元',
        },
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.suffix').exists()).toBe(true)
      expect(wrapper.find('.suffix').text()).toBe('元')
    })

    it('should not render prefix/suffix when not provided', () => {
      const wrapper = mount(StatCard, {
        props: defaultProps,
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.prefix').exists()).toBe(false)
      expect(wrapper.find('.suffix').exists()).toBe(false)
    })
  })

  describe('趋势显示', () => {
    it('should not show trend when not provided', () => {
      const wrapper = mount(StatCard, {
        props: defaultProps,
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.card-trend').exists()).toBe(false)
    })

    it('should show positive trend with trend-up class', () => {
      const wrapper = mount(StatCard, {
        props: {
          ...defaultProps,
          trend: 10.5,
        },
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.card-trend').exists()).toBe(true)
      expect(wrapper.find('.trend-value').classes()).toContain('trend-up')
      expect(wrapper.text()).toContain('10.50%')
    })

    it('should show negative trend with trend-down class', () => {
      const wrapper = mount(StatCard, {
        props: {
          ...defaultProps,
          trend: -5.25,
        },
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.trend-value').classes()).toContain('trend-down')
      expect(wrapper.text()).toContain('5.25%')
    })

    it('should show zero trend as trend-up', () => {
      const wrapper = mount(StatCard, {
        props: {
          ...defaultProps,
          trend: 0,
        },
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.trend-value').classes()).toContain('trend-up')
      expect(wrapper.text()).toContain('0.00%')
    })

    it('should show custom trend label', () => {
      const wrapper = mount(StatCard, {
        props: {
          ...defaultProps,
          trend: 10,
          trendLabel: '较上周',
        },
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.trend-label').text()).toBe('较上周')
    })

    it('should use default trend label "较昨日"', () => {
      const wrapper = mount(StatCard, {
        props: {
          ...defaultProps,
          trend: 10,
        },
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.trend-label').text()).toBe('较昨日')
    })
  })

  describe('值类型支持', () => {
    it('should accept string value', () => {
      const wrapper = mount(StatCard, {
        props: {
          title: '测试',
          value: '1,234.56',
        },
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.value').text()).toBe('1,234.56')
    })

    it('should accept number value', () => {
      const wrapper = mount(StatCard, {
        props: {
          title: '测试',
          value: 9999,
        },
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.value').text()).toBe('9999')
    })
  })

  describe('DOM 结构', () => {
    it('should have correct class structure', () => {
      const wrapper = mount(StatCard, {
        props: defaultProps,
        global: {
          stubs: {
            'el-icon': true,
            CaretTop: mockIcon,
            CaretBottom: mockIcon,
          },
        },
      })

      expect(wrapper.find('.stat-card').exists()).toBe(true)
      expect(wrapper.find('.card-content').exists()).toBe(true)
      expect(wrapper.find('.card-header').exists()).toBe(true)
      expect(wrapper.find('.card-title').exists()).toBe(true)
      expect(wrapper.find('.card-value').exists()).toBe(true)
      expect(wrapper.find('.value').exists()).toBe(true)
    })
  })
})
