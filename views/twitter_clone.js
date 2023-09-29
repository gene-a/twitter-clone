document.addEventListener('DOMContentLoaded', () => {
  ;(() => {
    /**
     * Helper function wrapping getElementById
     * @param {string} id
     * @returns
     */
    const getElById = (id) => {
      return document.getElementById(id)
    }
    const userCookieName = 'user'
    const apiEndPoint = '/api'
    const apiUserEndpoint = `${apiEndPoint}/user`
    const apiUserFollowerEndpoint = `${apiUserEndpoint}/follower`
    const apiFeedEndpoint = `${apiEndPoint}/feed`
    const apiTweetEndpoint = `${apiEndPoint}/tweet`
    const elFeed = getElById('appFeed')
    const elSidebar = getElById('appSidebar')
    const elTweetBox = getElById('appTweetBox')
    const elFeedHeader = getElById('feedHeader')
    const elPostTemplate = getElById('postTemplate')
    const elFeedContainer = getElById('feedContainer')
    const elDiscoverPane = getElById('discoverPane')
    const elDiscoverPaneTweet = getElById('suggestedUser')
    const elLightbox = getElById('lightbox')
    const elLoginBtn = getElById('loginButton')
    const elLogoutBtn = getElById('logoutButton')
    const elSuggestedTweetFollowButton = getElById('suggestedTweetFollowButton')
    const elSuggestedTweetAnotherButton = getElById(
      'suggestedTweetAnotherButton'
    )
    const elCreateTweetBtn = getElById('createTweetButton')
    const elFeedShowMoreBtn = getElById('feedShowMoreButton')
    const elSnackbar = getElById('snackbar')
    const elUsernameInput = getElById('username')
    const elPasswordInput = getElById('password')
    const elCreateTweetInput = getElById('userTweet')
    const feedPollSize = 3
    let currentFeedPage = 0

    /**
     * Toggle the visibility of the elemenet
     * @param {object} component
     */
    const toggleElementDisplay = (component) => {
      component.classList.toggle('hide')
      component.classList.toggle('show')
    }

    /**
     * Display the snackbar
     * @param {string} message
     */
    const toggleSnackbar = (message) => {
      elSnackbar.innerText = message
      elSnackbar.className = 'display'
      setTimeout(function () {
        elSnackbar.className = elSnackbar.className.replace('display', '')
      }, 3000)
    }

    /**
     * Helper function for updating the feed header
     * @param {string} message
     */
    const updateFeedHeader = (message) => {
      elFeedHeader.innerText = message
    }

    /**
     * Helper function for clearing the feed
     */
    const clearFeedAndDiscoverPane = () => {
      currentFeedPage = 0
      // Clear the pdiscover pane
      while (elDiscoverPaneTweet.firstChild) {
        elDiscoverPaneTweet.removeChild(elDiscoverPaneTweet.firstChild)
      }

      // Clear the feed
      while (elFeedContainer.firstChild) {
        elFeedContainer.removeChild(elFeedContainer.firstChild)
      }
    }

    /**
     * Helper function for toggling app visibility
     */
    const toggleAppVisibility = () => {
      toggleElementDisplay(elLightbox)
      toggleElementDisplay(elSidebar)
      toggleElementDisplay(elFeed)
      toggleElementDisplay(elTweetBox)
      toggleElementDisplay(elDiscoverPane)
    }

    /**
     * Handles the onclick event for the logout button
     */
    const onClickLogoutBtn = () => {
      if (getCookie(userCookieName)) {
        deleteCookie(userCookieName)
        toggleAppVisibility()
        // Reset the page feed always
        clearFeedAndDiscoverPane()
      }
    }

    /**
     * Handles the onclick event for the login button
     */
    const onClickLoginBtn = async () => {
      const requestBody = {
        username: elUsernameInput.value,
        password: elPasswordInput.value,
      }

      try {
        const response = await fetchPost(
          `${apiUserEndpoint}/login`,
          requestBody
        )

        if (!response.ok) {
          throw new Error('Network response was not ok')
        }

        const data = await response.json()

        if (data) {
          // Store user info into a cookie that lasts for 2 days, key is 'user'
          setCookie(userCookieName, JSON.stringify(data), 2)

          // Close the lightbox if user was authenticated
          // Display the rest of the app
          toggleAppVisibility()

          // Clear input
          elPasswordInput.value = ''
          elUsernameInput.value = ''

          // Display toast
          toggleSnackbar(`Hello! ${data.username}`)

          // Update header
          updateFeedHeader(`Home Feed - ${data.username}`)

          // Reset the page feed always
          clearFeedAndDiscoverPane()

          // Load the user's feed
          pollUserFeed()

          // Render discover pane
          pollSuggestedTweet()
        } else {
          // Create an error object and throw it
          throw new Error('Invalid username or password')
        }
      } catch (error) {
        // Handle any errors that occurred during the fetch
        toggleSnackbar(error.message)
        console.error('Post error:', error)
      }
    }

    /**
     * Handles the onclick event for the show more tweets button
     */
    const onClickShowMoreTweetsBtn = () => {
      pollUserFeed()
    }

    /**
     * Handles the onclick event for the create tweet button
     */
    const onClickCreateTweetBtn = async (e) => {
      e.preventDefault()
      const user = JSON.parse(getCookie(userCookieName))
      const requestBody = {
        user_id: user.id,
        tweet: elCreateTweetInput.value,
      }

      try {
        const response = await fetchPost(`${apiTweetEndpoint}`, requestBody)

        if (!response.ok) {
          throw new Error('Network response was not ok')
        }

        const data = await response.json()

        if (data) {
          toggleSnackbar('Tweet created!')
          // Create the new tweet as a feed item
          let post = createFeedItemFromTemplate(data)

          // Append the new tweet to the top
          elFeedContainer.prepend(post)

          // Clear the input field
          elCreateTweetInput.value = ''
        } else {
          // Create an error object and throw it
          throw new Error('Failed to create tweet')
        }
      } catch (error) {
        // Handle any errors that occurred during the fetch
        toggleSnackbar(error.message)
        console.error('Post error:', error)
      }
    }

    /**
     * Handles the onclick event for the get another tweet button
     */
    const onClickSuggestedTweetAnotherButton = async () => {
      if (!document.getElementById('suggestedUser').childElementCount) {
        // When there's no more suggested tweets buttons should be clickable but don't do anything
        toggleSnackbar('You have followed everyone! Wow!')
        return
      }
      pollSuggestedTweet()
      toggleSnackbar('Got a new suggested tweet!')
    }

    /**
     * Handles the onclick event for following a user in the suggested tweet page
     */
    const onClickSuggestedTweetFollowButton = async () => {
      if (!document.getElementById('suggestedUser').childElementCount) {
        // When there's no more suggested tweets buttons should be clickable but don't do anything
        toggleSnackbar('You have followed everyone! Wow!')
        return
      }

      try {
        const user = JSON.parse(getCookie(userCookieName))
        const suggestedUserId = document
          .querySelector('#suggestedUser .postUsername')
          .getAttribute('data-u-id')
        // The suggested user is the one we are following hence user_id being assigned to
        // Vice versa the current user is the follower hence follower_user_id
        const requestBody = {
          user_id: parseInt(suggestedUserId),
          follower_user_id: user.id,
        }
        const response = await fetchPut(
          `${apiUserFollowerEndpoint}`,
          requestBody
        )

        if (!response.ok) {
          throw new Error('Network response was not ok')
        }

        toggleSnackbar('Followed the user! Showing a new suggested tweet!')
        pollSuggestedTweet()
      } catch (error) {
        // Handle any errors that occurred during the fetch
        toggleSnackbar(error.message)
        console.error('Post error:', error)
      }
    }

    /**
     * Function for polling a suggested tweet from a user we are not following
     */
    const pollSuggestedTweet = async () => {
      const userJson = JSON.parse(getCookie(userCookieName))

      const requestParams = {
        id: userJson.id,
      }
      const response = await fetchGet(
        `${apiFeedEndpoint}/suggested?id=${requestParams.id}`,
        ''
      )

      if (!response.ok) {
        throw new Error('Network response was not ok')
      }

      const feedData = await response.json()

      if (feedData.id == 0) {
        while (elDiscoverPaneTweet.firstChild) {
          elDiscoverPaneTweet.removeChild(elDiscoverPaneTweet.firstChild)
        }

        // Handle the case when no suggested tweets are found
        toggleSnackbar(
          'No suggested tweets found, you might be following everyone. Wow!'
        )
      } else if (feedData) {
        // Loop through the data and create duplicates
        let post = createFeedItemFromTemplate(feedData)

        // Clear the pane then append the customized post to the discover pane
        while (elDiscoverPaneTweet.firstChild) {
          elDiscoverPaneTweet.removeChild(elDiscoverPaneTweet.firstChild)
        }
        elDiscoverPaneTweet.appendChild(post)
      } else {
        // When there is no suggested tweet data, we display to the user that everything's been loaded for them
        toggleSnackbar('You have followed everyone! Wow!')
      }
    }

    /**
     * Function for rendering the user feed
     */
    const pollUserFeed = async () => {
      const userJson = JSON.parse(getCookie(userCookieName))

      // Increment the current page, we begin with 0
      currentFeedPage += 1

      const requestParams = {
        id: userJson.id,
        page: currentFeedPage,
        pageSize: feedPollSize,
      }
      const response = await fetchGet(
        `${apiFeedEndpoint}?id=${requestParams.id}&page=${requestParams.page}&pageSize=${requestParams.pageSize}`,
        ''
      )

      if (!response.ok) {
        throw new Error('Network response was not ok')
      }

      const feedData = await response.json()

      if (feedData) {
        // Loop through the data and create duplicates
        feedData.forEach((data) => {
          let post = createFeedItemFromTemplate(data)

          // Append the customized post to the container
          elFeedContainer.appendChild(post)
        })
      } else {
        // When Feed data is empty we display to the user that everything's been loaded for them
        toggleSnackbar('No more tweets! You are all caught up!')
      }
    }

    /**
     * Utility function for creating a feed item from a template
     * @param {JSON} data
     */
    const createFeedItemFromTemplate = (data) => {
      // Clone the template
      const feedItem = elPostTemplate.cloneNode(true)

      // Clear the ID as the template should be the only one who has an ID
      feedItem.removeAttribute('id')

      // Remove the 'hide' class to make it visible
      feedItem.classList.toggle('hide')

      // Customize the content
      feedItem.querySelector('#postUsername').textContent = `@${data.username}`
      feedItem
        .querySelector('#postUsername')
        .setAttribute('data-u-id', data.user_id)
      feedItem.querySelector('#postContent').textContent = data.tweet

      // Add the class names for future selection work
      feedItem.querySelector('#postUsername').classList.toggle('postUsername')
      feedItem.querySelector('#postContent').classList.toggle('postContent')

      // Clear the IDs from the feed item, we only retain it during templating
      feedItem.querySelector('#postUsername').removeAttribute('id')
      feedItem.querySelector('#postContent').removeAttribute('id')

      return feedItem
    }

    /**
     * Helper function that wraps JS fetch function with the method of PUT
     * Sends a request body whether provided or not
     * @param {string} apiEndpoint
     * @param {object} requestBody
     * @returns Promise for api request
     */
    const fetchPut = async (apiEndpoint, requestBody) => {
      return fetch(apiEndpoint, {
        method: 'PUT',
        mode: 'same-origin',
        headers: {
          'Content-Type': 'application/json',
          'Accept-Encoding': 'gzip, deflate, br',
        },
        body: JSON.stringify(requestBody),
      })
    }

    /**
     * Helper function that wraps JS fetch function with the method of POST
     * Sends a request body whether provided or not
     * @param {string} apiEndpoint
     * @param {object} requestBody
     * @returns Promise for api request
     */
    const fetchPost = async (apiEndpoint, requestBody) => {
      return fetch(apiEndpoint, {
        method: 'POST',
        mode: 'same-origin',
        headers: {
          'Content-Type': 'application/json',
          'Accept-Encoding': 'gzip, deflate, br',
        },
        body: JSON.stringify(requestBody),
      })
    }

    /**
     * Helper function that wraps JS fetch function with the method of GET
     * Sends a requestbody whether provided or not
     * @param {string} apiEndpoint
     * @param {object} requestBody
     * @returns Promise for api request
     */
    const fetchGet = async (apiEndpoint, requestBody) => {
      return fetch(apiEndpoint, {
        method: 'GET',
        mode: 'same-origin',
        headers: {
          'Accept-Encoding': 'gzip, deflate, br',
        },
      })
    }

    /**
     * Util function for storing cookies
     * @param {string} name
     * @param {object} value
     * @param {number} days
     */
    const setCookie = (name, value, days) => {
      let expires = ''
      if (name && days) {
        let date = new Date()
        date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000)
        expires = '; expires=' + date.toUTCString()
      }

      document.cookie = name + '=' + value + expires + '; path=/'
    }

    /**
     * Util function for getting cookies
     * @param {string} name
     * @returns
     */
    const getCookie = (name) => {
      if (!name) return null

      let nameEQ = name + '='
      let cookies = document.cookie.split(';')
      for (let i = 0; i < cookies.length; i++) {
        var cookie = cookies[i]
        while (cookie.charAt(0) == ' ') {
          cookie = cookie.substring(1, cookie.length)
        }
        if (cookie.indexOf(nameEQ) == 0) {
          return cookie.substring(nameEQ.length, cookie.length)
        }
      }
      return null
    }

    /**
     * Util function for deleting a cookie
     * @param {string} name
     */
    const deleteCookie = (name) => {
      if (!name) return
      // To delete a cookie, we set its expiration date to a date in the past
      document.cookie =
        name + '=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/'
    }

    /**
     * Init function for the module
     */
    const init = function () {
      // Event listener for the login button
      elLoginBtn.addEventListener('click', onClickLoginBtn)
      // Event listener for the logout button
      elLogoutBtn.addEventListener('click', onClickLogoutBtn)
      // Event listener for the load more tweets button
      elFeedShowMoreBtn.addEventListener('click', onClickShowMoreTweetsBtn)
      // Event listener for the get another tweet suggestion button
      elSuggestedTweetAnotherButton.addEventListener(
        'click',
        onClickSuggestedTweetAnotherButton
      )
      // Event listener for the follow suggested user button
      elSuggestedTweetFollowButton.addEventListener(
        'click',
        onClickSuggestedTweetFollowButton
      )
      // Event listener for the create tweet button
      elCreateTweetBtn.addEventListener('click', (e) => {
        onClickCreateTweetBtn(e)
      })

      // We check if there are cookies available, if yes we log the user in automatically
      if (document.cookie) {
        const user = JSON.parse(getCookie(userCookieName))
        // Close the lightbox since we found the user cookie
        // Display the rest of the app
        toggleAppVisibility()

        // Display toast
        toggleSnackbar(`Welcome back! ${user.username}`)

        // Update header
        updateFeedHeader(`Home Feed - ${user.username}`)

        // Reset the page feed always
        clearFeedAndDiscoverPane()

        // Load the user's feed
        pollUserFeed()

        // Render discover pane
        pollSuggestedTweet()
      }
    }
    return {
      init: init,
    }
  })().init()
})
